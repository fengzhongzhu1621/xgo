package connpool

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/fengzhongzhu1621/xgo/network/dial"
)

type dialFunc = func(ctx context.Context) (net.Conn, error)

var _ IPool = (*pool)(nil)

// pool connection pool factory, maintains connection pools corresponding to all addresses,
// and connection pool option information.
// pool 维护一个 sync.Map 作为连接池，
// key 为<network, address, protocol>编码，value 为与目标地址建立的连接构成的 ConnectionPool, 其内部以一个链表维护空闲连接。
// 在短连接模式中，transport 层会在 rpc 调用后关闭连接，而在连接池模式中，会把使用完的连接放回连接池，以待下次需要时取出。
type pool struct {
	opts            *Options  // 连接池配置
	connectionPools *sync.Map // 连接池映射，key 为<network, address, protocol>编码，value 为与目标地址建立的连接构成的 ConnectionPool, 其内部以一个链表维护空闲连接。
}

// DefaultConnectionPool is the default connection pool, replaceable.
var DefaultConnectionPool = NewConnectionPool()

// NewConnectionPool creates a connection pool.
// NewConnectionPool 创建一个连接池，支持传入 Option 修改参数，不传则使用默认值初始化。
// Dial 是默认的创建连接方式，每个 ConnectionPool 会根据自己的 GetOptions 生成 DialOptions, 来建立对应目标的连接。
func NewConnectionPool(opt ...Option) IPool {
	// Default value, tentative, need to debug to determine the specific value.
	opts := &Options{
		MaxIdle:         defaultMaxIdle,         // 最大空闲连接数量
		IdleTimeout:     defaultIdleTimeout,     // 空闲连接超时时间
		DialTimeout:     defaultDialTimeout,     // 建立连接超时时间
		PoolIdleTimeout: defaultPoolIdleTimeout, // 连接池空闲超时时间
		Dial:            dial.Dial,              // 发起连接请求的方法
	}
	for _, o := range opt {
		o(opts)
	}
	return &pool{
		opts:            opts,
		connectionPools: new(sync.Map),
	}
}

// Get is used to get the connection from the connection pool.
// 获取一个连接，仅暴露 Get 接口，确保连接池状态不会因用户的误操作被破坏。
func (p *pool) Get(network string, address string, opts GetOptions) (net.Conn, error) {
	// 获取拨号操作的上下文
	ctx, cancel := opts.getDialCtx(p.opts.DialTimeout)
	if cancel != nil {
		defer cancel()
	}

	// 生成连接池节点的唯一标识键
	key := getNodeKey(network, address, opts.Protocol)

	// 从连接池中获取连接池
	if v, ok := p.connectionPools.Load(key); ok {
		// 如果连接池已存在，使用现有的连接池获取连接
		return v.(*ConnectionPool).Get(ctx)
	}

	// 如果连接池不存在，创建新的连接池
	newPool := &ConnectionPool{
		Dial:               p.getDialFunc(network, address, opts),
		MinIdle:            p.opts.MinIdle,
		MaxIdle:            p.opts.MaxIdle,
		MaxActive:          p.opts.MaxActive,
		Wait:               p.opts.Wait,
		MaxConnLifetime:    p.opts.MaxConnLifetime,
		IdleTimeout:        p.opts.IdleTimeout,
		framerBuilder:      opts.FramerBuilder,
		customReader:       opts.CustomReader,
		forceClosed:        p.opts.ForceClose,
		PushIdleConnToTail: p.opts.PushIdleConnToTail,
		onCloseFunc:        func() { p.connectionPools.Delete(key) },
		poolIdleTimeout:    p.opts.PoolIdleTimeout,
	}

	// 如果设置了最大活跃连接数，创建令牌通道用于连接数控制
	if newPool.MaxActive > 0 {
		newPool.tokenChan = make(chan struct{}, p.opts.MaxActive)
	}

	// 设置连接检查器，优先使用自定义检查器
	newPool.checker = newPool.defaultChecker
	if p.opts.Checker != nil {
		newPool.checker = p.opts.Checker // 使用自定义检查器
	}

	// 避免在初始化期间并发写入连接池映射的问题
	v, ok := p.connectionPools.LoadOrStore(key, newPool)
	if !ok {
		// 如果是新创建的连接池，进行初始化和注册检查器
		newPool.RegisterChecker(defaultCheckInterval, newPool.checker)
		newPool.keepMinIdles() // 保持最小空闲连接数
		return newPool.Get(ctx)
	}
	// 如果连接池已存在，使用现有的连接池获取连接
	return v.(*ConnectionPool).Get(ctx)
}

// getDialFunc 创建拨号函数，用于建立网络连接
// 参数:
//
//	network: 网络协议类型
//	address: 目标地址
//	opts: 获取连接的配置选项
//
// 返回:
//
//	dialFunc: 拨号函数
func (p *pool) getDialFunc(network string, address string, opts GetOptions) dialFunc {
	// 配置拨号选项
	dialOpts := &dial.DialOptions{
		Network:       network,            // 网络协议
		Address:       address,            // 目标地址
		LocalAddr:     opts.LocalAddr,     // 本地绑定地址
		CACertFile:    opts.CACertFile,    // CA证书文件
		TLSCertFile:   opts.TLSCertFile,   // 客户端证书文件
		TLSKeyFile:    opts.TLSKeyFile,    // 客户端私钥文件
		TLSServerName: opts.TLSServerName, // 服务器名称
		IdleTimeout:   p.opts.IdleTimeout, // 空闲超时时间
	}

	// 返回实际的拨号函数
	return func(ctx context.Context) (net.Conn, error) {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err() // 上下文已取消，返回错误
		default:
		}

		// 获取上下文的截止时间
		d, ok := ctx.Deadline()
		if !ok {
			return nil, ErrNoDeadline // 没有设置截止时间，返回错误
		}

		// 复制拨号选项并设置超时时间
		opts := *dialOpts
		opts.Timeout = time.Until(d) // 计算剩余超时时间

		// 调用实际的拨号方法
		return p.opts.Dial(&opts)
	}
}

// getNodeKey 生成连接池节点的唯一标识键
// 参数:
//
//	network: 网络协议类型（如 "tcp", "udp"）
//	address: 目标地址（格式为 "host:port"）
//	protocol: 协议名称
//
// 返回:
//
//	string: 格式为 "network_address_protocol" 的唯一键
func getNodeKey(network, address, protocol string) string {
	const underline = "_" // 分隔符
	var key strings.Builder

	// 预分配足够的缓冲区空间以提高性能
	// 计算总长度：network长度 + address长度 + protocol长度 + 2个下划线
	key.Grow(len(network) + len(address) + len(protocol) + 2)

	// 构建键值：network_address_protocol
	key.WriteString(network)
	key.WriteString(underline)
	key.WriteString(address)
	key.WriteString(underline)
	key.WriteString(protocol)

	return key.String()
}
