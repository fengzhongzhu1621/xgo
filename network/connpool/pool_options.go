package connpool

import (
	"time"

	"github.com/fengzhongzhu1621/xgo/network/dial"
)

// Option is the Options helper.
type Option func(*Options)

// Options indicates pool configuration.
type Options struct {
	// 最小空闲连接数，由连接池后台周期性补入，0 表示不补入
	MinIdle int // Initialize the number of connections, ready for the next io.
	// 最大空闲连接数量，0 不做限制
	MaxIdle int // Maximum number of idle connections, 0 means no idle.
	// 用户可用连接的最大并发数，0 不做限制
	MaxActive int // Maximum number of active connections, 0 means no limit.
	// Whether to wait when the maximum number of active connections is reached.
	// 是否等待
	Wait bool
	// 空闲连接超时时间，0 不做限制
	IdleTimeout time.Duration // idle connection timeout.
	// 连接的最大生命周期，0 不做限制
	MaxConnLifetime time.Duration // Maximum lifetime of the connection.
	// 建立连接超时时间
	DialTimeout time.Duration // Connection establishment timeout.
	// 用户使用连接后是否强制关闭，适用于流式场景
	ForceClose bool
	// 发起连接请求的方法
	Dial    dial.DialFunc
	Checker IHealthChecker
	// 放回连接池的方式，默认采用 LIFO 获取空闲连接
	PushIdleConnToTail bool // connection to ip will be push tail when ConnectionPool.put method is called
	// 连接池空闲超时时间， 0 表示不做检查
	PoolIdleTimeout time.Duration // ConnectionPool idle timeout
}

// WithMinIdle returns an Option which sets the number of initialized connections.
func WithMinIdle(n int) Option {
	return func(o *Options) {
		o.MinIdle = n
	}
}

// WithMaxIdle returns an Option which sets the maximum number of idle connections.
func WithMaxIdle(m int) Option {
	return func(o *Options) {
		o.MaxIdle = m
	}
}

// WithMaxActive returns an Option which sets the maximum number of active connections.
func WithMaxActive(s int) Option {
	return func(o *Options) {
		o.MaxActive = s
	}
}

// WithWait returns an Option which sets whether to wait when the number of connections reaches the limit.
func WithWait(w bool) Option {
	return func(o *Options) {
		o.Wait = w
	}
}

// WithIdleTimeout returns an Option which sets the idle connection time, after which it may be closed.
func WithIdleTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.IdleTimeout = t
	}
}

// WithMaxConnLifetime returns an Option which sets the maximum lifetime of
// the connection, after which it may be closed.
func WithMaxConnLifetime(t time.Duration) Option {
	return func(o *Options) {
		o.MaxConnLifetime = t
	}
}

// WithDialTimeout returns an Option which sets the default timeout time for
// the connection pool to establish a connection.
func WithDialTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.DialTimeout = t
	}
}

// WithForceClose returns an Option which sets whether to force the connection to be closed.
func WithForceClose(f bool) Option {
	return func(o *Options) {
		o.ForceClose = f
	}
}

// WithDialFunc returns an Option which sets dial function.
func WithDialFunc(d dial.DialFunc) Option {
	return func(o *Options) {
		o.Dial = d
	}
}

// WithHealthChecker returns an Option which sets health checker.
func WithHealthChecker(c IHealthChecker) Option {
	return func(o *Options) {
		o.Checker = c
	}
}

// WithPushIdleConnToTail returns an Option which sets PushIdleConnToTail flag.
func WithPushIdleConnToTail(c bool) Option {
	return func(o *Options) {
		o.PushIdleConnToTail = c
	}
}

// WithPoolIdleTimeout returns an Option which sets pool idle timeout.
// after the timeout, ConnectionPool resource may be cleaned up.
func WithPoolIdleTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.PoolIdleTimeout = t
	}
}
