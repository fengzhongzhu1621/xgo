package server_transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/fengzhongzhu1621/xgo/network/reuseport"
	"github.com/fengzhongzhu1621/xgo/network/transport/handler"
	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/panjf2000/ants/v2"
)

const (
	// EnvGraceRestart is the flag of graceful restart.
	EnvGraceRestart = "XGO_RPC_IS_GRACEFUL"

	// EnvGraceFirstFd is the fd of graceful first listener.
	EnvGraceFirstFd = "XGO_RPC_GRACEFUL_1ST_LISTENFD"

	// EnvGraceRestartFdNum is the number of fd for graceful restart.
	EnvGraceRestartFdNum = "XGO_RPC_GRACEFUL_LISTENFD_NUM"

	// EnvGraceRestartPPID is the PPID of graceful restart.
	EnvGraceRestartPPID = "XGO_RPC_GRACEFUL_PPID"
)

var (
	errUnSupportedListenerType = errors.New("not supported listener type")
	errUnSupportedNetworkType  = errors.New("not supported network type")
	errFileIsNotSocket         = errors.New("file is not a socket")
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// conn is the struct of connection which is established when server receive a client connecting
// request.
type conn struct {
	ctx         context.Context
	cancelCtx   context.CancelFunc
	idleTimeout time.Duration
	lastVisited time.Time
	handler     handler.IHandler
}

func (c *conn) handle(ctx context.Context, req []byte) ([]byte, error) {
	return c.handler.Handle(ctx, req)
}

func (c *conn) handleClose(ctx context.Context) error {
	if closeHandler, ok := c.handler.(handler.ICloseHandler); ok {
		return closeHandler.HandleClose(ctx)
	}
	return nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var _ IServerTransport = (*serverTransport)(nil)

// serverTransport 服务器传输的实现细节，支持TCP或UDP
type serverTransport struct {
	addrToConn map[string]*tcpconn             // 地址到TCP连接的映射
	m          *sync.RWMutex                   // 读写锁，保护连接映射的并发安全
	opts       *options.ServerTransportOptions // 服务器传输配置选项
}

// DefaultServerTransport ServerStreamTransport的默认实现
var DefaultServerTransport = NewServerTransport(options.WithReusePort(true))

// NewServerTransport 创建新的IServerTransport实例
func NewServerTransport(opt ...options.ServerTransportOption) IServerTransport {
	r := newServerTransport(opt...)
	return &r
}

// newServerTransport 创建新的serverTransport实例
func newServerTransport(opt ...options.ServerTransportOption) serverTransport {
	// 使用默认选项
	opts := options.DefaultServerTransportOptions()
	for _, o := range opt {
		o(opts) // 应用用户提供的选项
	}
	addrToConn := make(map[string]*tcpconn) // 初始化连接映射
	return serverTransport{addrToConn: addrToConn, m: &sync.RWMutex{}, opts: opts}
}

// ListenAndServe 开始监听，失败时返回错误
func (s *serverTransport) ListenAndServe(ctx context.Context, opts ...options.ListenServeOption) error {
	lsopts := &options.ListenServeOptions{}
	for _, opt := range opts {
		opt(lsopts) // 应用监听服务选项
	}

	if lsopts.Listener != nil {
		return s.listenAndServeStream(ctx, lsopts) // 使用自定义监听器
	}
	// 支持同时监听TCP和UDP
	networks := strings.Split(lsopts.Network, ",")
	for _, network := range networks {
		lsopts.Network = network
		switch lsopts.Network {
		case "tcp", "tcp4", "tcp6", "unix":
			if err := s.listenAndServeStream(ctx, lsopts); err != nil {
				return err
			}
		case "udp", "udp4", "udp6":
			if err := s.listenAndServePacket(ctx, lsopts); err != nil {
				return err
			}
		default:
			return fmt.Errorf("server transport: not support network type %s", lsopts.Network)
		}
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ---------------------------------stream server-----------------------------------------//

var (
	// listenersMap records the listeners in use in the current process.
	listenersMap = &sync.Map{}
	// inheritedListenersMap record the listeners inherited from the parent process.
	// A key(host:port) may have multiple listener fds.
	inheritedListenersMap = &sync.Map{}
	// once controls fds passed from parent process to construct listeners.
	once sync.Once
)

// listenAndServeStream 启动流式监听，失败时返回错误
func (s *serverTransport) listenAndServeStream(ctx context.Context, opts *options.ListenServeOptions) error {
	if opts.FramerBuilder == nil {
		return errors.New("tcp transport FramerBuilder empty")
	}

	// 获取TCP监听器
	ln, err := s.getTCPListener(opts)
	if err != nil {
		return fmt.Errorf("get tcp listener err: %w", err)
	}

	// 必须保存原始TCP监听器（而不是TLS监听器）以确保热重启时可以成功检索底层文件描述符
	listenersMap.Store(ln, struct{}{})

	// 可能升级为TLS监听器
	ln, err = mayLiftToTLSListener(ln, opts)
	if err != nil {
		return fmt.Errorf("may lift to tls listener err: %w", err)
	}

	// 异步启动 TCP 流式服务
	go s.serveStream(ctx, ln, opts)
	return nil
}

// newConn 创建新的连接对象
// 参数:
//   - ctx: 上下文，用于取消和超时控制
//   - opts: 监听服务选项，包含连接配置信息
//
// 返回值: 新创建的连接对象指针
func (s *serverTransport) newConn(ctx context.Context, opts *options.ListenServeOptions) *conn {
	// 从监听服务选项中获取空闲超时时间
	idleTimeout := opts.IdleTimeout

	// 检查服务器传输选项中是否设置了空闲超时时间
	// 如果设置了，则优先使用服务器传输选项中的超时时间
	if s.opts.IdleTimeout > 0 {
		idleTimeout = s.opts.IdleTimeout
	}

	// 创建并返回新的连接对象
	return &conn{
		ctx:         ctx,          // 设置连接上下文
		handler:     opts.Handler, // 设置业务处理器
		idleTimeout: idleTimeout,  // 设置空闲超时时间
	}
}

// getTCPListener 获取/创建 TCP/Unix监听器
func (s *serverTransport) getTCPListener(opts *options.ListenServeOptions) (listener net.Listener, err error) {
	listener = opts.Listener // 使用自定义监听器

	if listener != nil {
		return listener, nil // 直接返回自定义监听器
	}

	// 检查优雅重启环境变量，环境变量不存在 v 为空，ok 为 false
	v, _ := os.LookupEnv(EnvGraceRestart)
	ok, _ := strconv.ParseBool(v)
	if ok {
		// 查找从父进程传递的监听器（优雅重启）
		pln, err := getPassedListener(opts.Network, opts.Address)
		if err != nil {
			return nil, err
		}

		listener, ok := pln.(net.Listener)
		if !ok {
			return nil, errors.New("invalid net.Listener")
		}
		return listener, nil // 返回继承的监听器
	}

	// 端口复用：为了加速IO，内核将IO ReadReady事件分发给线程
	if s.opts.ReusePort && opts.Network != "unix" {
		listener, err = reuseport.Listen(opts.Network, opts.Address) // 使用端口复用监听
		if err != nil {
			return nil, fmt.Errorf("%s reuseport error:%v", opts.Network, err)
		}
	} else {
		listener, err = net.Listen(opts.Network, opts.Address) // 标准网络监听
		if err != nil {
			return nil, err
		}
	}

	return listener, nil // 返回创建的监听器
}

// serveStream 启动 TCP 流式服务
func (s *serverTransport) serveStream(ctx context.Context, ln net.Listener, opts *options.ListenServeOptions) error {
	var once sync.Once
	closeListener := func() { ln.Close() } // 关闭监听器的函数
	defer once.Do(closeListener)           // 确保函数退出时关闭监听器

	// 创建goroutine监听关闭信号
	// 一旦Server.Close()被调用，TCP监听器应立即关闭并不再接受新连接
	go func() {
		select {
		case <-ctx.Done():
			// ctx.Done会执行以下两个操作：
			// 1. 停止监听
			// 2. 取消所有当前已建立的连接
			// 而opts.StopListening只会停止监听
		case <-opts.StopListening:
		}
		logging.Tracef("recv server close event") // 记录服务器关闭事件
		once.Do(closeListener)                    // 执行一次关闭监听器操作
	}()

	// 启动TCP服务
	return s.serveTCP(ctx, ln, opts)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ---------------------------------packet server-----------------------------------------//
// listenAndServePacket 启动数据包监听，失败时返回错误
func (s *serverTransport) listenAndServePacket(ctx context.Context, opts *options.ListenServeOptions) error {
	pool := createUDPRoutinePool(opts.Routines) // 创建UDP协程池
	// 端口复用：为了加速IO，内核将IO ReadReady事件分发给线程
	if s.opts.ReusePort {
		reuseport.ListenerBacklogMaxSize = 4096
		cores := runtime.NumCPU()
		for i := 0; i < cores; i++ {
			udpconn, err := s.getUDPListener(opts) // 获取UDP监听器
			if err != nil {
				return err
			}
			listenersMap.Store(udpconn, struct{}{})

			go s.servePacket(ctx, udpconn, pool, opts) // 启动数据包服务
		}
	} else {
		udpconn, err := s.getUDPListener(opts)
		if err != nil {
			return err
		}
		listenersMap.Store(udpconn, struct{}{})

		go s.servePacket(ctx, udpconn, pool, opts)
	}
	return nil
}

// getUDPListener gets UDP listener.
func (s *serverTransport) getUDPListener(opts *options.ListenServeOptions) (udpConn net.PacketConn, err error) {
	v, _ := os.LookupEnv(EnvGraceRestart)
	ok, _ := strconv.ParseBool(v)
	if ok {
		// Find the passed listener.
		ln, err := getPassedListener(opts.Network, opts.Address)
		if err != nil {
			return nil, err
		}
		listener, ok := ln.(net.PacketConn)
		if !ok {
			return nil, errors.New("invalid net.PacketConn")
		}
		return listener, nil
	}

	if s.opts.ReusePort {
		udpConn, err = reuseport.ListenPacket(opts.Network, opts.Address)
		if err != nil {
			return nil, fmt.Errorf("udp reuseport error:%v", err)
		}
	} else {
		udpConn, err = net.ListenPacket(opts.Network, opts.Address)
		if err != nil {
			return nil, fmt.Errorf("udp listen error:%v", err)
		}
	}

	return udpConn, nil
}

func (s *serverTransport) servePacket(ctx context.Context, rwc net.PacketConn, pool *ants.PoolWithFunc,
	opts *options.ListenServeOptions) error {
	switch rwc := rwc.(type) {
	case *net.UDPConn:
		return s.serveUDP(ctx, rwc, pool, opts)
	default:
		return errors.New("transport not support PacketConn impl")
	}
}
