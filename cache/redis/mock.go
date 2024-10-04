package redis

import (
	"fmt"
	"time"

	redisClient "github.com/fengzhongzhu1621/xgo/db/redisx/client"
	"github.com/go-redis/cache/v8"
)

// NewMockCache will create a cache for mock
// name: 缓存的名称，将用于生成缓存键的前缀
// expiration: 缓存项的默认过期时间
// returns: 返回一个新的 Cache 结构体实例
func NewMockCache(name string, expiration time.Duration) *Cache {
	// 创建一个新的 Redis 客户端实例，该实例连接到由 miniredis 模拟的 Redis 服务器
	cli := redisClient.NewTestRedisClient()

	// key format = xgo:{cache_name}:{real_key}
	keyPrefix := fmt.Sprintf("xgo:%s", name)

	// 创建了一个新的缓存编解码器（codec）。cache.New 函数接受一个 cache.Options 结构体，其中指定了 Redis 客户端
	codec := cache.New(&cache.Options{
		Redis: cli,
	})

	// 返回一个新的 Cache 结构体实例
	return &Cache{
		name:              name,       // 缓存的名称
		keyPrefix:         keyPrefix,  // 缓存键的前缀
		codec:             codec,      // 缓存编解码器
		cli:               cli,        // Redis 客户端实例
		defaultExpiration: expiration, // 缓存项的默认过期时间
	}
}
