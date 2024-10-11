package backoff

import (
	"sync/atomic"
	"time"
)

// 自定义锁重试策略的接口
type IRetryStrategy interface {
	// NextBackoff returns the next backoff duration 获取下一次重试的等待时间
	NextBackoff() time.Duration
}

////////////////////////////////////////////////////////////////////////////////////////////

type linearBackoff time.Duration

// NextBackoff 下一次重试前的时间间隔，实现了一个恒定的重试间隔。
func (r linearBackoff) NextBackoff() time.Duration {
	return time.Duration(r)
}

func NewLinearBackoff(backoff time.Duration) IRetryStrategy {
	return linearBackoff(backoff)
}

////////////////////////////////////////////////////////////////////////////////////////////

// NoRetry 其等待时间为0，意味着不进行重试，只尝试一次。
func NoRetry() IRetryStrategy {
	return linearBackoff(0)
}

// //////////////////////////////////////////////////////////////////////////////////////////
// limitedRetry 定义重试策略，以及两个int64类型的字段cnt和max，分别用于记录当前的重试次数和最大重试次数。
type limitedRetry struct {
	s   IRetryStrategy // 重试策略
	cnt int64          // 当前的重试次数
	max int64          // 最大重试次数
}

// LimitRetry limits the number of retries to max attempts.
func NewLimitRetry(s IRetryStrategy, max int) IRetryStrategy {
	return &limitedRetry{s: s, max: int64(max)}
}

// NextBackoff 获取下一次重试的等待时间
func (r *limitedRetry) NextBackoff() time.Duration {
	if atomic.LoadInt64(&r.cnt) >= r.max {
		return 0
	}
	// 使用原子操作增加重试次数计数器
	atomic.AddInt64(&r.cnt, 1)
	return r.s.NextBackoff()
}
