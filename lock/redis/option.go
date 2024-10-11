package option

import (
	"github.com/fengzhongzhu1621/xgo/collections/backoff"
)

// Options 用于配置分布式锁的选项
type Options struct {
	// 获取下一次重试的等待时间的接口
	RetryStrategy backoff.IRetryStrategy

	// 用于附加到锁令牌上
	Metadata string
}

// GetMetadata 用于获取 Options 结构体中的 Metadata 字段的值
func (o *Options) GetMetadata() string {
	if o != nil {
		return o.Metadata
	}
	return ""
}

// GetRetryStrategy 获取下一次重试的等待时间策略
func (o *Options) GetRetryStrategy() backoff.IRetryStrategy {
	if o != nil && o.RetryStrategy != nil {
		return o.RetryStrategy
	}
	// 默认是不等待
	return backoff.NoRetry()
}
