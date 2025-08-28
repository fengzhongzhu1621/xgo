package backoff

import (
	"time"
)

// 获取下一次重试的等待时间的接口
type IRetryStrategy interface {
	// NextBackoff returns the next backoff duration 获取下一次重试的等待时间
	NextBackoff() time.Duration
}

// NoRetry 其等待时间为0，意味着不进行重试，只尝试一次。
func NoRetry() IRetryStrategy {
	return linearBackoff(0)
}
