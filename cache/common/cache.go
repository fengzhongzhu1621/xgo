package common

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// NewTTLCache 创建一个新的内存缓存实例，并返回该实例的指针 create cache with expiration and cleanup interval,
// if cleanupInterval is 0, will use DefaultCleanupInterval
//
// - expiration 缓存项的过期时间。过了这个时间，缓存项将被视为过期并可能被清理。
// - cleanupInterval 缓存清理的间隔时间。在这个间隔内，缓存系统会检查并清理过期的缓存项。
func NewTTLCache(expiration time.Duration, cleanupInterval time.Duration) *gocache.Cache {
	// 默认 5 分钟清除一次缓存
	if cleanupInterval == 0 {
		cleanupInterval = DefaultCleanupInterval
	}

	return gocache.New(expiration, cleanupInterval)
}
