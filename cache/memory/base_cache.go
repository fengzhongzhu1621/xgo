package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/cache"
	"github.com/fengzhongzhu1621/xgo/cache/memory/backend"
	"golang.org/x/sync/singleflight"
)

const EmptyCacheExpiration = 5 * time.Second

// EmptyCache is a placeholder for the missing key
type EmptyCache struct {
	err error
}

// BaseCache is a cache which retrieves data from the backend and stores it in the cache.
type BaseCache struct {
	// 缓存后台，支持多个第三方缓存库
	backend backend.Backend

	// 是否禁用缓存
	disabled bool
	// 禁用缓存时根据 key 获取 value 的函数
	retrieveFunc             RetrieveFunc
	g                        singleflight.Group
	withEmptyCache           bool
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
		// if retrieve fail from retrieveFunc
		// 如果缓存的值存在但是value 的值是 EmptyCache，则返回错误
		if emptyCache, isEmptyCache := value.(EmptyCache); isEmptyCache {
			return nil, emptyCache.err
		}
		return value, nil
	}

	// 3. if not exists in cache, retrieve it
	// 缓存不存在时获取值，并缓存执行结果
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
			// 记录错误信息
			c.backend.Set(key, EmptyCache{err: err}, c.emptyCacheExpireDuration)
		}
		return nil, err
	}

	// 4. set value to cache, use default expiration
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
		return defaultZeroTime, fmt.Errorf("not a string value. key=%s, value=%v(%T)", k.Key(), value, value)
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

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Option func(*BaseCache)

func WithNoCache() Option {
	return func(cache *BaseCache) {
		cache.disabled = true
	}
}

// WithEmptyCache will set the key EmptyCache  if retrieve fail from retrieveFunc
func WithEmptyCache(timeout time.Duration) Option {
	return func(cache *BaseCache) {
		if timeout == 0 {
			timeout = EmptyCacheExpiration
		}
		cache.withEmptyCache = true
		cache.emptyCacheExpireDuration = timeout
	}
}

func NewBaseCache(retrieveFunc RetrieveFunc, backend backend.Backend, options ...Option) Cache {
	c := &BaseCache{
		backend:      backend,
		retrieveFunc: retrieveFunc,
	}
	// 自定义参数
	for _, o := range options {
		o(c)
	}
	return c
}