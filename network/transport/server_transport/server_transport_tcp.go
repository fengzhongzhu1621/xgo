package server_transport

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/fengzhongzhu1621/xgo/buildin/buffer"
	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/collections/ring/writev"
	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/fengzhongzhu1621/xgo/network/ip"
	"github.com/fengzhongzhu1621/xgo/network/transport/frame"
	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/fengzhongzhu1621/xgo/xerror"
	"github.com/panjf2000/ants/v2"
)

// tcpconn is the connection which is established when server accept a client connecting request.
type tcpconn struct {
	*conn
	rwc         net.Conn
	fr          codec.IFramer
	localAddr   net.Addr
	remoteAddr  net.Addr
	serverAsync bool
	writev      bool
	copyFrame   bool
	closeOnce   sync.Once
	st          *serverTransport
	pool        *ants.PoolWithFunc
	buffer      *writev.Buffer
	closeNotify chan struct{}
}

// close closes socket and cleans up.
func (c *tcpconn) close() {
	c.closeOnce.Do(func() {
		// Send error msg to handler.
		ctx, msg := codec.WithNewMessage(context.Background())
		msg.WithLocalAddr(c.localAddr)
		msg.WithRemoteAddr(c.remoteAddr)
		e := &xerror.Error{
			Type: xerror.ErrorTypeFramework,
			Code: xerror.RetServerSystemErr,
			Desc: "trpc",
			Msg:  "Server connection closed",
		}
		msg.WithServerRspErr(e)
		// The connection closing message is handed over to handler.
		if err := c.conn.handleClose(ctx); err != nil {
			logging.Trace("transport: notify connection close failed", err)
		}
		// Notify to stop writev sending goroutine.
		if c.writev {
			close(c.closeNotify)
		}

		// Remove cache in server stream transport.
		key := ip.AddrToKey(c.localAddr, c.remoteAddr)
		c.st.m.Lock()
		delete(c.st.addrToConn, key)
		c.st.m.Unlock()

		// Finally, close the socket connection.
		c.rwc.Close()
	})
}

// write encapsulates tcp conn write.
func (c *tcpconn) write(p []byte) (int, error) {
	if c.writev {
		return c.buffer.Write(p)
	}
	return c.rwc.Write(p)
}

func (c *tcpconn) serve() {
	defer c.close()
	for {
		// Check if upstream has closed.
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		if c.idleTimeout > 0 {
			now := time.Now()
			// SetReadDeadline has poor performance, so, update timeout every 5 seconds.
			if now.Sub(c.lastVisited) > 5*time.Second {
				c.lastVisited = now
				err := c.rwc.SetReadDeadline(now.Add(c.idleTimeout))
				if err != nil {
					logging.Trace("transport: tcpconn SetReadDeadline fail ", err)
					return
				}
			}
		}

		// 读取一个完整的包
		req, err := c.fr.ReadFrame()
		if err != nil {
			if err == io.EOF {
				return
			}
			// Server closes the connection if client sends no package in last idle timeout.
			if e, ok := err.(net.Error); ok && e.Timeout() {
				return
			}
			logging.Trace("transport: tcpconn serve ReadFrame fail ", err)
			return
		}

		// if framer is not concurrent safe, copy the data to avoid over writing.
		if c.copyFrame {
			reqCopy := make([]byte, len(req))
			copy(reqCopy, req)
			req = reqCopy
		}

		// 处理业务逻辑
		c.handle(req)
	}
}

// handle 处理业务逻辑
// 如果开启了异步处理，则将处理参数放入协程池中，否则直接调用handleSyncWithErr函数处理
func (c *tcpconn) handle(req []byte) {
	if !c.serverAsync || c.pool == nil {
		c.handleSync(req)
		return
	}

	// Using sync.pool to dispatch package processing goroutine parameters can reduce a memory
	// allocation and slightly promote performance.
	args := handleParamPool.Get().(*handleParam)
	args.req = req
	args.c = c
	args.start = time.Now()
	if err := c.pool.Invoke(args); err != nil {
		logging.Trace("transport: tcpconn serve routine pool put job queue fail ", err)
		c.handleSyncWithErr(req, xerror.ErrServerRoutinePoolBusy)
	}
}

// handleSyncWithErr 同步处理业务逻辑
func (c *tcpconn) handleSync(req []byte) {
	c.handleSyncWithErr(req, nil)
}

// handleSyncWithErr 处理业务逻辑
func (c *tcpconn) handleSyncWithErr(req []byte, e error) {
	ctx, msg := codec.WithNewMessage(context.Background())
	defer codec.PutBackMessage(msg)

	msg.WithServerRspErr(e)
	// Record local addr and remote addr to context.
	msg.WithLocalAddr(c.localAddr)
	msg.WithRemoteAddr(c.remoteAddr)

	// 处理业务逻辑
	rsp, err := c.conn.handle(ctx, req)

	if err != nil {
		if err != xerror.ErrServerNoResponse {
			logging.Trace("transport: tcpconn serve handle fail ", err)
			c.close()
			return
		}
		// On stream RPC, server does not need to write rsp, just returns.
		return
	}

	{
		// common RPC write rsp.
		_, err = c.write(rsp)
	}

	if err != nil {
		logging.Trace("transport: tcpconn write fail ", err)
		c.close()
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *serverTransport) serveTCP(ctx context.Context, ln net.Listener, opts *options.ListenServeOptions) error {
	// 创建协程池（如果开启了 ServerAsync 选项）
	var pool *ants.PoolWithFunc
	if opts.ServerAsync {
		// 创建协程池，并设置最大协程数为opts.Routines
		// 协程池的任务函数为handleParam.run，该函数会调用handleParam.c.handleSyncWithErr函数
		pool = createRoutinePool(opts.Routines)
	}

	// tempDelay 用于临时错误的重试延迟，采用指数退避策略
	for tempDelay := time.Duration(0); ; {
		// 接受新的TCP连接，这是一个阻塞调用，每个连接都会启动一个tc.serve()单独处理
		rwc, err := ln.Accept()
		if err != nil {
			// 处理接受连接时出现的错误
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				// 如果是临时性错误，采用指数退避策略计算下一次延迟时间，并在阻塞等待下一次连接
				tempDelay = doTempDelay(tempDelay)
				continue
			}

			select {
			case <-ctx.Done():
				// 如果上下文被取消（如服务器关闭或重启），直接返回错误退出监听循环
				return err
			default:
				// 处理其他类型的错误
				const accept, closeError = "accept", "use of closed network connection"
				const msg = "the server transport, listening on %s, encountered an error: %+v; this error was handled" +
					" gracefully by the framework to prevent abnormal termination, serving as a reference for" +
					" investigating acceptance errors that can't be filtered by the Temporary interface"

				// 检查是否为连接关闭错误（通常发生在服务器正常关闭时）
				if e, ok := err.(*net.OpError); ok && e.Op == accept && strings.Contains(e.Err.Error(), closeError) {
					logging.Infof("listener with address %s is closed", ln.Addr())
					return err
				}
				// 记录其他无法通过Temporary接口过滤的错误
				logging.Errorf(msg, ln.Addr(), err)
				continue
			}
		}

		// 重置临时错误延迟计数器
		tempDelay = 0

		// 对TCP连接进行优化配置
		if tcpConn, ok := rwc.(*net.TCPConn); ok {
			// 启用TCP KeepAlive机制，检测连接是否存活
			if err := tcpConn.SetKeepAlive(true); err != nil {
				logging.Tracef("tcp conn set keepalive error:%v", err)
			}
			// 设置KeepAlive探测间隔
			if s.opts.KeepAlivePeriod > 0 {
				if err := tcpConn.SetKeepAlivePeriod(s.opts.KeepAlivePeriod); err != nil {
					logging.Tracef("tcp conn set keepalive period error:%v", err)
				}
			}
		}

		// 创建TCP连接包装对象，包含连接管理所需的所有信息
		tc := &tcpconn{
			conn:        s.newConn(ctx, opts),                          // 创建新的连接对象
			rwc:         rwc,                                           // 原始网络连接
			fr:          opts.FramerBuilder.New(buffer.NewReader(rwc)), // 帧解析器，用于处理数据帧
			remoteAddr:  rwc.RemoteAddr(),                              // 客户端地址
			localAddr:   rwc.LocalAddr(),                               // 服务器地址
			serverAsync: opts.ServerAsync,                              // 是否启用异步处理模式
			writev:      opts.Writev,                                   // 是否使用writev系统调用进行向量写操作
			st:          s,                                             // 指向serverTransport的引用
			pool:        pool,                                          // 协程池引用（如果启用）
		}

		// 如果启用writev优化，初始化写缓冲区
		if tc.writev {
			tc.buffer = writev.NewBuffer()          // 创建写缓冲区
			tc.closeNotify = make(chan struct{}, 1) // 关闭通知通道
			tc.buffer.Start(tc.rwc, tc.closeNotify) // 启动写缓冲区处理例程
		}

		// 检查是否需要对数据帧进行拷贝，避免数据覆盖问题
		// 判断依据：配置选项、异步模式是否启用、帧解析器是否线程安全
		tc.copyFrame = frame.ShouldCopy(opts.CopyFrame, tc.serverAsync, codec.IsSafeFramer(tc.fr))

		// 生成连接的唯一标识键，用于连接管理
		key := ip.AddrToKey(tc.localAddr, tc.remoteAddr)
		s.m.Lock()
		s.addrToConn[key] = tc // 将连接注册到连接映射表中
		s.m.Unlock()

		// 启动新的goroutine处理该连接，实现并发处理
		go tc.serve()
	}
}

// doTempDelay 实现了指数退避策略，用于处理临时性错误的重试延迟
// 参数 tempDelay: 当前的延迟时间，如果是第一次重试则为0
// 返回值: 计算出的新延迟时间，用于下一次重试
func doTempDelay(tempDelay time.Duration) time.Duration {
	// 判断是否为第一次重试
	if tempDelay == 0 {
		// 第一次重试，设置初始延迟为5毫秒
		tempDelay = 5 * time.Millisecond
	} else {
		// 非第一次重试，将延迟时间翻倍（指数增长核心）
		tempDelay *= 2
	}

	// 设置最大延迟上限为1秒，防止延迟时间无限增长
	if max := 1 * time.Second; tempDelay > max {
		tempDelay = max
	}

	// 让当前goroutine休眠指定的延迟时间
	time.Sleep(tempDelay)

	// 返回新的延迟时间，供下一次重试使用
	return tempDelay
}
