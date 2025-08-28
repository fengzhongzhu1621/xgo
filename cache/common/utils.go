package common

import (
	"math/rand"
	"time"
)

// NewRandomDuration 生成一个最大 1 秒的随机偏移
func NewRandomDuration(seconds int) RandomExtraExpirationDurationFunc {
	return func() time.Duration {
		return time.Duration(rand.Intn(seconds*1000)) * time.Millisecond
	}
}
