package backend

import (
	"time"

	cache "github.com/fengzhongzhu1621/xgo/cache/common"
	gocache "github.com/patrickmn/go-cache"
)

// MemoryBackend is the backend for memory cache
type MemoryBackend struct {
	name  string
	cache *gocache.Cache

	// 缓存项的过期时间。过了这个时间，缓存项将被视为过期并可能被清理。
	defaultExpiration time.Duration
	// 生成随机过期时间偏移量的函数，用于打散过期时间
	randomExtraExpirationFunc cache.RandomExtraExpirationDurationFunc
}

// NewMemoryBackend create memory backend
// - name: the name of the backend
// - expiration: the expiration duration of the cache 缓存项的过期时间。过了这个时间，缓存项将被视为过期并可能被清理。
// - randomExtraExpirationFunc: the function generate extra expiration duration 生成随机过期时间函数
func NewMemoryBackend(
	name string,
	expiration time.Duration,
	randomExtraExpirationFunc cache.RandomExtraExpirationDurationFunc,
) *MemoryBackend {
	cleanupInterval := expiration + (5 * time.Minute)

	return &MemoryBackend{
		name:                      name,
		cache:                     cache.NewTTLCache(expiration, cleanupInterval),
		defaultExpiration:         expiration,
		randomExtraExpirationFunc: randomExtraExpirationFunc,
	}
}

// Set sets value to cache with key and expiration
func (c *MemoryBackend) Set(key string, value interface{}, duration time.Duration) {
	if duration == time.Duration(0) {
		duration = c.defaultExpiration
	}

	// 过期时间 + 随机偏移量，用于打散过期时间
	if c.randomExtraExpirationFunc != nil {
		duration += c.randomExtraExpirationFunc()
	}

	c.cache.Set(key, value, duration)
}

// Get gets value by key from the cache
func (c *MemoryBackend) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

// Delete deletes value by key from the cache
func (c *MemoryBackend) Delete(key string) error {
	c.cache.Delete(key)
	return nil
}
