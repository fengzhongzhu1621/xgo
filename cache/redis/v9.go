package redis

import (
	"context"
	"fmt"
	"time"

	redisClient "github.com/fengzhongzhu1621/xgo/db/redis/client"
	"github.com/go-redis/cache/v9"
)

type CacheV9 struct {
	// 缓存的名称，用于标识和区分不同的缓存实例
	name string
	// 缓存键的前缀，用于在 Redis 中为所有相关的缓存键添加统一的前缀，以避免键名冲突
	keyPrefix string
	// 编解码器，用于序列化和反序列化缓存的数据。通常是一个实现了缓存数据编码和解码逻辑的对象
	cache *cache.Cache
	// 默认的缓存过期时间，当设置缓存但没有明确指定过期时间时，将使用此默认值
	ttl time.Duration
}

func NewCacheV9(name string, ttl time.Duration) *CacheV9 {
	keyPrefix := fmt.Sprintf("xgo:09:%s", name)

	cli := redisClient.GetDefaultRedisV9Client()

	c := cache.New(&cache.Options{
		Redis:      cli,
		LocalCache: nil,
	})
	return &CacheV9{
		name:      name,
		keyPrefix: keyPrefix,
		cache:     c,
		ttl:       ttl,
	}
}

// genKey 生成一个带有前缀的缓存键
func (c *CacheV9) genKey(key string) string {
	return c.keyPrefix + ":" + key
}

func (c *CacheV9) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	// 0 表示默认值
	if ttl == time.Duration(0) {
		ttl = c.ttl
	}

	return c.cache.Set(&cache.Item{Ctx: ctx, Key: c.genKey(key), Value: value, TTL: ttl})
}

func (c *CacheV9) Exists(ctx context.Context, key string) bool {
	// 和 v8 的不同在于对 Exists 返回值格式的处理
	return c.cache.Exists(ctx, c.genKey(key))
}

func (c *CacheV9) Get(ctx context.Context, key string, value any) error {
	return c.cache.Get(ctx, c.genKey(key), value)
}

func (c *CacheV9) Delete(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, c.genKey(key))
}
