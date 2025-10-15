package multiplexed

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	defaultConnNumberPerHost = 2
	defaultSendQueueSize     = 1024
	defaultDialTimeout       = time.Second
)

// Multiplexed 表示多路复用连接池管理器
type Multiplexed struct {
	mu sync.RWMutex // 读写锁，保护连接映射的并发安全
	// 三级连接映射结构:
	// key(ip:port)
	//   => value(*Connections)         <-- 同一ip:port的多个具体连接集合
	//     => (*Connection)             <-- 特定ip:port的单个具体连接
	//       => [](*VirtualConnection)  <-- 在特定具体连接上复用的多个虚拟连接
	concreteConns *sync.Map    // 具体连接映射表，key为节点键，value为Connections
	opts          *PoolOptions // 连接池配置选项
}

// DefaultMultiplexedPool 是默认的多路复用连接池实现
var DefaultMultiplexedPool = New()

// New 创建新的多路复用连接池实例
// 参数:
//
//	opt: 可选的连接池配置选项
//
// 返回:
//
//	*Multiplexed: 新创建的多路复用连接池
func New(opt ...PoolOption) *Multiplexed {
	opts := &PoolOptions{
		connectNumberPerHost: defaultConnNumberPerHost, // 每个地址的连接数量
		sendQueueSize:        defaultSendQueueSize,     // 每个连接的发送队列长度
		dialTimeout:          defaultDialTimeout,       // 连接超时时间，默认1秒
	}
	for _, o := range opt {
		o(opts)
	}

	// 最大空闲连接数不能小于预分配连接数
	if opts.maxIdleConnsPerHost != 0 && opts.maxIdleConnsPerHost < opts.connectNumberPerHost {
		opts.maxIdleConnsPerHost = opts.connectNumberPerHost
	}
	return &Multiplexed{
		concreteConns: new(sync.Map),
		opts:          opts,
	}
}

// GetMuxConn 获取到指定网络地址的多路复用连接
// 参数:
//
//	ctx: 上下文，用于控制超时和取消
//	network: 网络协议类型（如"tcp"、"udp"）
//	address: 目标地址（格式为"host:port"）
//	opts: 连接获取配置选项
//
// 返回:
//
//	IMuxConn: 多路复用连接接口
//	error: 获取连接过程中发生的错误
func (p *Multiplexed) GetMuxConn(
	ctx context.Context,
	network string,
	address string,
	opts GetOptions,
) (IMuxConn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := opts.update(network, address); err != nil {
		return nil, err
	}

	// 获取虚拟连接
	return p.get(ctx, &opts)
}

// get 内部方法，执行获取虚拟连接的三步流程
// 参数:
//
//	ctx: 上下文
//	opts: 连接获取配置选项
//
// 返回:
//
//	*VirtualConnection: 虚拟连接
//	error: 获取过程中发生的错误
func (p *Multiplexed) get(ctx context.Context, opts *GetOptions) (*VirtualConnection, error) {
	// 步骤1: 节点键(ip:port) => 具体连接集合
	value, ok := p.concreteConns.Load(opts.nodeKey)
	if !ok {
		p.initPoolForNode(opts)
		value, ok = p.concreteConns.Load(opts.nodeKey)
		if !ok {
			return nil, ErrInitPoolFail
		}
	}
	conns, ok := value.(*Connections)
	if !ok {
		return nil, fmt.Errorf("%w, expected: *Connections, actual: %T", ErrAssertFail, value)
	}

	// 步骤2: 具体连接集合 => 单个具体连接
	conn, err := conns.pickSingleConcrete(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf(
			"multiplexed pick single concreate connection with node key %s err: %w", opts.nodeKey, err)
	}

	// 步骤3: 单个具体连接 => 虚拟连接
	return conn.newVirConn(ctx, opts.VID), nil
}

// initPoolForNode 为指定节点初始化连接池
// 参数:
//
//	opts: 连接获取配置选项
func (p *Multiplexed) initPoolForNode(opts *GetOptions) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 再次检查，以防另一个goroutine在我们之前初始化了连接池
	if _, ok := p.concreteConns.Load(opts.nodeKey); ok {
		return
	}
	// 创建新的具体连接集合
	p.concreteConns.Store(opts.nodeKey, p.newConcreteConnections(opts))
}

// newConcreteConnections 创建新的具体连接集合
// 参数:
//
//	opts: 连接获取配置选项
//
// 返回:
//
//	*Connections: 新创建的具体连接集合
func (p *Multiplexed) newConcreteConnections(opts *GetOptions) *Connections {
	conns := &Connections{
		nodeKey: opts.nodeKey,
		opts:    p.opts,
		conns:   make([]*Connection, 0, p.opts.connectNumberPerHost),
		maxIdle: p.opts.maxIdleConnsPerHost,
		destructor: func() {
			p.concreteConns.Delete(opts.nodeKey)
		},
	}

	conns.initialize(opts)

	return conns
}
