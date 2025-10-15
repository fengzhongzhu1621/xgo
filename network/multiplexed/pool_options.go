package multiplexed

import "time"

// PoolOptions 表示连接池的一些配置设置
type PoolOptions struct {
	connectNumberPerHost int           // 每个地址的连接数量
	sendQueueSize        int           // 每个连接的发送队列长度
	dropFull             bool          // 队列满时是否丢弃请求
	dialTimeout          time.Duration // 连接超时时间，默认1秒
	maxVirConnsPerConn   int           // 每个真实连接的最大虚拟连接数，0表示无限制
	maxIdleConnsPerHost  int           // 每个对端ip:port的最大空闲连接数
}

// PoolOption 是配置选项的辅助类型
type PoolOption func(*PoolOptions)

// WithConnectNumber 设置连接池中每个对端的连接数量
// 参数:
//   number: 连接数量
func WithConnectNumber(number int) PoolOption {
	return func(opts *PoolOptions) {
		opts.connectNumberPerHost = number
	}
}

// WithQueueSize 设置连接池中每个连接的发送队列长度
// 参数:
//   n: 队列长度
func WithQueueSize(n int) PoolOption {
	return func(opts *PoolOptions) {
		opts.sendQueueSize = n
	}
}

// WithDropFull 设置当队列满时是否丢弃请求
// 参数:
//   drop: true表示丢弃，false表示阻塞等待
func WithDropFull(drop bool) PoolOption {
	return func(opts *PoolOptions) {
		opts.dropFull = drop
	}
}

// WithDialTimeout 设置连接超时时间
// 参数:
//   d: 超时时间
func WithDialTimeout(d time.Duration) PoolOption {
	return func(opts *PoolOptions) {
		opts.dialTimeout = d
	}
}

// WithMaxVirConnsPerConn 设置每个真实连接的最大虚拟连接数
// 参数:
//   n: 虚拟连接数，0表示无限制
func WithMaxVirConnsPerConn(n int) PoolOption {
	return func(opts *PoolOptions) {
		opts.maxVirConnsPerConn = n
	}
}

// WithMaxIdleConnsPerHost 设置每个对端ip:port的最大空闲连接数
// 此值不应小于ConnectNumber，通常与MaxVirConnsPerConn在流式场景中配合使用
// 用于动态调整连接数量，仅在MaxVirConnsPerConn设置时生效，0表示无限制
// 参数:
//   n: 最大空闲连接数
func WithMaxIdleConnsPerHost(n int) PoolOption {
	return func(opts *PoolOptions) {
		opts.maxIdleConnsPerHost = n
	}
}
