package memory

import (
	"context"
	"fmt"
	"time"

	cache "github.com/fengzhongzhu1621/xgo/cache/common"
	"github.com/fengzhongzhu1621/xgo/cache/memory/backend"
	"golang.org/x/sync/singleflight"
)

// BaseCache is a cache which retrieves data from the backend and stores it in the cache.
type BaseCache struct {
	// 缓存的后端存储，可以是任何实现了 backend.Backend 接口的第三方缓存库
	backend backend.Backend

	// 一个布尔值，用于指示缓存是否被禁用。
	disabled bool
	// 禁用缓存时根据 key 获取 value 的函数
	retrieveFunc cache.RetrieveFunc
	// 用于防止缓存击穿，即多个请求同时访问不存在的缓存项时，只执行一次数据获取操作。
	g singleflight.Group
	// 用于指示是否启用空缓存机制
	withEmptyCache bool
	// 空缓存的过期时间
	emptyCacheExpireDuration time.Duration
}

// Exists returns true if the cache has a value for the given key.
func (c *BaseCache) Exists(ctx context.Context, key cache.Key) bool {
	k := key.Key()
	_, ok := c.backend.Get(k)
	return ok
}

// Get will get the key from cache, if missing, will call the retrieveFunc to get the data, add to cache, then return
func (c *BaseCache) Get(ctx context.Context, key cache.Key) (interface{}, error) {
	// 1. if cache is disabled, fetch and return
	if c.disabled {
		// 缓存禁用后，使用指定方法获取 value
		value, err := c.retrieveFunc(ctx, key)
		if err != nil {
			return nil, err
		}
		return value, nil
	}

	k := key.Key()

	// 2. get from cache
	value, ok := c.backend.Get(k)
	if ok {
		// 处理无效的缓存值, 如果从缓存中获取到的值是一个特殊的 EmptyCache 类型（表示该键曾经存在但现在已被删除或过期），则返回该 EmptyConfig 实例中包含的错误信息。
		// EmptyCache 可能是一个自定义的类型，用于表示缓存中曾经存在但现已无效的条目。这种设计允许缓存系统在不删除键的情况下标记某些条目为无效，
		// 这在某些复杂的缓存策略中可能很有用。
		if emptyCache, isEmptyCache := value.(cache.EmptyCache); isEmptyCache {
			return nil, emptyCache.Err
		}
		return value, nil
	}

	// 3. if not exists in cache, retrieve it
	// 如果键既不在缓存中，也不是无效的 EmptyCache 类型，则调用 doRetrieve 方法来从数据源（可能是数据库、API 等）获取该键对应的值，并将其添加到缓存中。最后返回新获取到的值。
	return c.doRetrieve(ctx, key)
}

// doRetrieve 缓存不存在时获取值，并缓存执行结果 will retrieve the real data from database, redis, apis, etc.
func (c *BaseCache) doRetrieve(ctx context.Context, k cache.Key) (interface{}, error) {
	key := k.Key()

	// 3.2 fetch
	// Do(key string, fn func() (interface{}, error)) (interface{}, error)：执行给定的函数 fn，并返回其结果。
	// 如果有多个 goroutine 同时调用此方法并使用相同的 key，则只有一个 goroutine 会执行 fn，
	// 其他 goroutine 会等待结果并共享相同的返回值。
	value, err, _ := c.g.Do(key, func() (interface{}, error) {
		return c.retrieveFunc(ctx, k)
	})

	if err != nil {
		if c.withEmptyCache {
			// ! if error, cache it too, make it short enough(5s)
			// 检查是否启用了空缓存机制（withEmptyCache）。如果启用了，它会将一个特殊的 EmptyCache 对象（包含错误信息）存入缓存，并设置一个较短的过期时间（通常是 5 秒）。
			// 这样做可以避免频繁地尝试获取同一个不存在的数据，同时给系统一些时间来恢复。
			c.backend.Set(key, cache.EmptyCache{Err: err}, c.emptyCacheExpireDuration)
		}
		return nil, err
	}

	// 4. set value to cache, use default expiration
	// 如果成功获取到了数据，方法会将这个值存入缓存。这里使用 0 作为过期时间，表示使用缓存的默认过期策略。
	c.backend.Set(key, value, 0)

	return value, nil
}

// Set will set key-value into cache.
func (c *BaseCache) Set(ctx context.Context, key cache.Key, data interface{}) {
	k := key.Key()
	c.backend.Set(k, data, 0)
}

// Delete deletes the value from the cache for the given key.
func (c *BaseCache) Delete(ctx context.Context, key cache.Key) error {
	k := key.Key()
	return c.backend.Delete(k)
}

// DirectGet will get key from cache, without calling the retrieveFunc
func (c *BaseCache) DirectGet(ctx context.Context, key cache.Key) (interface{}, bool) {
	k := key.Key()
	return c.backend.Get(k)
}

// Disabled returns true if the cache is disabled.
func (c *BaseCache) Disabled() bool {
	return c.disabled
}

// GetString returns a string representation of the value for the given key.
// will error if the type is not a string.
func (c *BaseCache) GetString(ctx context.Context, k cache.Key) (string, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return "", err
	}

	v, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("not a string value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetBool returns a bool representation of the value for the given key.
// will error if the type is not a bool.
func (c *BaseCache) GetBool(ctx context.Context, k cache.Key) (bool, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return false, err
	}

	v, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("not a string value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

var defaultZeroTime = time.Time{}

// GetTime returns a time representation of the value for the given key.
// will error if the type is not an time.Time.
func (c *BaseCache) GetTime(ctx context.Context, k cache.Key) (time.Time, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return defaultZeroTime, err
	}

	v, ok := value.(time.Time)
	if !ok {
		return defaultZeroTime, fmt.Errorf(
			"not a string value. key=%s, value=%v(%T)",
			k.Key(),
			value,
			value,
		)
	}
	return v, nil
}

// GetInt returns an int representation of the value for the given key.
// will error if the type is not an int.
func (c *BaseCache) GetInt(ctx context.Context, k cache.Key) (int, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("not a int value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetInt8 returns an int8 representation of the value for the given key.
// will error if the type is not an int8.
func (c *BaseCache) GetInt8(ctx context.Context, k cache.Key) (int8, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(int8)
	if !ok {
		return 0, fmt.Errorf("not a int8 value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetInt16 returns an int16 representation of the value for the given key.
// will error if the type is not an int16.
func (c *BaseCache) GetInt16(ctx context.Context, k cache.Key) (int16, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(int16)
	if !ok {
		return 0, fmt.Errorf("not a int16 value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetInt32 returns an int32 representation of the value for the given key.
// will error if the type is not an int32.
func (c *BaseCache) GetInt32(ctx context.Context, k cache.Key) (int32, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(int32)
	if !ok {
		return 0, fmt.Errorf("not a int32 value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetInt64 returns an int64 representation of the value for the given key.
// will error if the type is not an int64.
func (c *BaseCache) GetInt64(ctx context.Context, k cache.Key) (int64, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("not a int64 value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetUint returns an uint representation of the value for the given key.
// will error if the type is not an uint.
func (c *BaseCache) GetUint(ctx context.Context, k cache.Key) (uint, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(uint)
	if !ok {
		return 0, fmt.Errorf("not a uint value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetUint8 returns an uint8 representation of the value for the given key.
// will error if the type is not an uint8.
func (c *BaseCache) GetUint8(ctx context.Context, k cache.Key) (uint8, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(uint8)
	if !ok {
		return 0, fmt.Errorf("not a uint8 value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetUint16 returns an uint16 representation of the value for the given key.
// will error if the type is not an uint16.
func (c *BaseCache) GetUint16(ctx context.Context, k cache.Key) (uint16, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(uint16)
	if !ok {
		return 0, fmt.Errorf("not a uint16 value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetUint32 returns an uint32 representation of the value for the given key.
// will error if the type is not an uint32.
func (c *BaseCache) GetUint32(ctx context.Context, k cache.Key) (uint32, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(uint32)
	if !ok {
		return 0, fmt.Errorf("not a uint32 value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetUint64 returns an uint64 representation of the value for the given key.
// will error if the type is not an uint64.
func (c *BaseCache) GetUint64(ctx context.Context, k cache.Key) (uint64, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(uint64)
	if !ok {
		return 0, fmt.Errorf("not a uint64 value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetFloat32 returns a float32 representation of the value for the given key.
// will error if the type is not a float32.
func (c *BaseCache) GetFloat32(ctx context.Context, k cache.Key) (float32, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(float32)
	if !ok {
		return 0, fmt.Errorf("not a float32 value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// GetFloat64 returns a float64 representation of the value for the given key.
// will error if the type is not a float64.
func (c *BaseCache) GetFloat64(ctx context.Context, k cache.Key) (float64, error) {
	value, err := c.Get(ctx, k)
	if err != nil {
		return 0, err
	}

	v, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("not a float64 value. key=%s, value=%v(%T)", k.Key(), value, value)
	}
	return v, nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Option func(*BaseCache)

func WithNoCache() Option {
	return func(cache *BaseCache) {
		cache.disabled = true
	}
}

// WithEmptyCache will set the key EmptyCache  if retrieve fail from retrieveFunc
func WithEmptyCache(timeout time.Duration) Option {
	return func(baseCache *BaseCache) {
		if timeout == 0 {
			timeout = cache.EmptyCacheExpiration
		}
		baseCache.withEmptyCache = true
		baseCache.emptyCacheExpireDuration = timeout
	}
}

func NewBaseCache(
	disabled bool,
	retrieveFunc cache.RetrieveFunc,
	backend backend.Backend,
	options ...Option,
) Cache {
	c := &BaseCache{
		backend:      backend,
		disabled:     disabled,
		retrieveFunc: retrieveFunc,
	}
	// 自定义参数
	for _, o := range options {
		o(c)
	}
	return c
}
