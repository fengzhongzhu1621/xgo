package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	gopkgcache "github.com/fengzhongzhu1621/xgo/cache/common"
	redisClient "github.com/fengzhongzhu1621/xgo/db/redisx/client"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/sync/singleflight"
)

type RetrieveFunc func(key gopkgcache.Key) (interface{}, error)

const (
	// while the go-redis/cache upgrade maybe not compatible with the previous version.
	// e.g. the object set by v7 can't read by v8
	// https://github.com/go-redis/cache/issues/52
	// NOTE: important!!! if upgrade the go-redis/cache version, should change the version

	// CacheVersion is loop in 00->99->00 => make sure will not conflict with previous version
	// 代码中提到go-redis/cache库的版本升级可能导致不兼容问题，需要注意版本管理
	CacheVersion = "00"

	// 批量 key 的最大数量
	PipelineSizeThreshold = 100
)

// Cache is a cache implements
// 定义了一个缓存实例，包含名称、键前缀、编解码器、Redis客户端和默认过期时间
type Cache struct {
	// 缓存的名称，用于标识和区分不同的缓存实例
	name string
	// 缓存键的前缀，用于在 Redis 中为所有相关的缓存键添加统一的前缀，以避免键名冲突
	keyPrefix string
	// 编解码器，用于序列化和反序列化缓存的数据。通常是一个实现了缓存数据编码和解码逻辑的对象
	codec *cache.Cache
	// Redis 客户端实例，用于与 Redis 服务器进行通信，执行实际的缓存读写操作
	cli *redis.Client
	// 默认的缓存过期时间，当设置缓存但没有明确指定过期时间时，将使用此默认值
	defaultExpiration time.Duration
	// 用于防止缓存击穿攻击。它可以确保在高并发情况下，对于同一缓存键的多次访问只会触发一次实际的数据加载操作。
	G singleflight.Group
}

// NewCache create a cache instance
func NewCache(name string, expiration time.Duration) *Cache {
	cli := redisClient.GetDefaultRedisClient()

	// key format = xgo:{version}:{cache_name}:{real_key}
	keyPrefix := fmt.Sprintf("xgo:%s:%s", CacheVersion, name)

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

// genKey 生成一个带有前缀的缓存键
func (c *Cache) genKey(key string) string {
	return c.keyPrefix + ":" + key
}

// copyTo 将数据从一个对象复制到另一个对象，使用msgpack进行序列化和反序列化
func (c *Cache) copyTo(source interface{}, dest interface{}) error {
	b, err := msgpack.Marshal(source)
	if err != nil {
		return err
	}

	err = msgpack.Unmarshal(b, dest)
	return err
}

// Set 将键值对设置到缓存中
func (c *Cache) Set(key gopkgcache.Key, value interface{}, duration time.Duration) error {
	// 检查过期时间
	if duration == time.Duration(0) {
		duration = c.defaultExpiration
	}
	// 生成完整的缓存键
	k := c.genKey(key.Key())
	// 设置缓存项
	return c.codec.Set(&cache.Item{
		Key:   k,
		Value: value,
		TTL:   duration,
	})
}

// Get execute `get`
func (c *Cache) Get(key gopkgcache.Key, value interface{}) error {
	k := c.genKey(key.Key())
	return c.codec.Get(context.TODO(), k, value)
}

// Exists execute `exists`
func (c *Cache) Exists(key gopkgcache.Key) bool {
	k := c.genKey(key.Key())

	count, err := c.cli.Exists(context.TODO(), k).Result()

	return err == nil && count == 1
}

// GetInto will retrieve the data from cache and unmarshal into the obj
// 从缓存中获取数据并将其反序列化到指定的对象中。如果缓存中没有数据，则会调用提供的 retrieveFunc 函数来获取数据，并将其存入缓存
// key: 要查询的缓存键
// obj: 目标对象，数据将被反序列化到这个对象中
// retrieveFunc: 用于在缓存未命中时从数据源获取数据
func (c *Cache) GetInto(key gopkgcache.Key, obj interface{}, retrieveFunc RetrieveFunc) (err error) {
	// 1. get from cache, hit, return
	// 尝试从缓存中获取数据:
	err = c.Get(key, obj)
	if err == nil {
		return nil
	}

	// 2. if missing
	// 2.1 check the guard
	// 2.2 do retrieve
	// 如果缓存未命中，调用 retrieveFunc 从数据源获取数据
	// 确保在多个并发请求同时尝试获取同一个不存在的键时，只会执行一次数据获取操作。
	data, err, _ := c.G.Do(key.Key(), func() (interface{}, error) {
		return retrieveFunc(key)
	})
	// 2.3 do retrieve fail, make guard and return
	if err != nil {
		// if retrieve fail, should wait for few seconds for the missing-retrieve
		// c.makeGuard(key)
		return
	}

	// 3. set to cache
	// 将获取到的数据存入缓存
	errNotImportant := c.Set(key, data, 0)
	if errNotImportant != nil {
		log.Errorf("set to redis fail, key=%s, err=%s", key.Key(), errNotImportant)
	}

	// 注意: 基础类型无法通过 *obj = value 直接赋值，因此需要通过 copyTo 方法进行反序列化赋值。
	// 所以利用从缓存再次反序列化给对应指针赋值(相当于底层msgpack.unmarshal帮做了转换再次反序列化给对应指针赋值
	return c.copyTo(data, obj)
}

// Delete execute `del`
func (c *Cache) Delete(key gopkgcache.Key) (err error) {
	k := c.genKey(key.Key())

	ctx := context.TODO()

	_, err = c.cli.Del(ctx, k).Result()
	return err
}

// Expire execute `expire`
func (c *Cache) Expire(key gopkgcache.Key, duration time.Duration) error {
	if duration == time.Duration(0) {
		duration = c.defaultExpiration
	}

	k := c.genKey(key.Key())
	_, err := c.cli.Expire(context.TODO(), k, duration).Result()
	return err
}

// BatchDelete execute `del` with pipeline 批量删除缓存中的多个键
func (c *Cache) BatchDelete(keys []gopkgcache.Key) error {
	// 生成完整的缓存键列表
	newKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		newKeys = append(newKeys, c.genKey(key.Key()))
	}
	ctx := context.TODO()

	var err error
	if len(newKeys) < PipelineSizeThreshold {
		// 直接使用 Del 方法批量删除这些键
		_, err = c.cli.Del(ctx, newKeys...).Result()
	} else {
		// 使用 Redis 的管道（pipeline）功能来批量删除这些键。这样可以提高删除操作的效率，减少网络往返次数
		pipe := c.cli.Pipeline()

		for _, key := range newKeys {
			pipe.Del(ctx, key)
		}

		_, err = pipe.Exec(ctx)
	}
	return err
}

// BatchExpireWithTx execute `expire` with tx pipeline
func (c *Cache) BatchExpireWithTx(keys []gopkgcache.Key, expiration time.Duration) error {
	pipe := c.cli.TxPipeline()
	ctx := context.TODO()

	for _, k := range keys {
		key := c.genKey(k.Key())
		pipe.Expire(ctx, key, expiration)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// KV is a key-value pair
type KV struct {
	Key   string
	Value string
}

// BatchGet execute `get` with pipeline 批量获取多个键的值
func (c *Cache) BatchGet(keys []gopkgcache.Key) (map[gopkgcache.Key]string, error) {
	pipe := c.cli.Pipeline()

	ctx := context.TODO()

	// 准备命令
	cmds := map[gopkgcache.Key]*redis.StringCmd{}
	for _, k := range keys {
		key := c.genKey(k.Key())
		cmd := pipe.Get(ctx, key)

		cmds[k] = cmd
	}

	// 执行管道中的所有命令，并获取执行结果
	_, err := pipe.Exec(ctx)
	// 当批量操作, 里面有个key不存在, err = redis.Nil; 但是不应该影响其他存在的key的获取
	// Nil reply returned by Redis when key does not exist.
	// redis.Nil 错误表示某个键在 Redis 中不存在，但这不应该影响其他键的获取
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// 收集结果
	values := make(map[gopkgcache.Key]string, len(cmds))
	for hkf, cmd := range cmds {
		// maybe err or key missing
		// only return the HashKeyField who get value success from redis
		val, err := cmd.Result()
		if err != nil {
			continue
		} else {
			// 只有成功获取到的键值对才会被包含在最终的返回结果中
			values[hkf] = val
		}
	}

	return values, nil
}

// BatchSetWithTx execute `set` with tx pipeline
func (c *Cache) BatchSetWithTx(kvs []KV, expiration time.Duration) error {
	if expiration == time.Duration(0) {
		expiration = c.defaultExpiration
	}

	// tx, all success or all fail
	// 创建一个事务管道
	// 使用事务管道可以确保所有设置操作要么全部成功，要么全部失败，从而保持数据的一致性。
	pipe := c.cli.TxPipeline()

	ctx := context.TODO()

	for _, kv := range kvs {
		key := c.genKey(kv.Key)
		pipe.Set(ctx, key, kv.Value, expiration)
	}

	// 如果所有操作成功，则返回 nil；如果有错误，则返回相应的错误
	_, err := pipe.Exec(ctx)
	return err
}

// ZData is a sorted-set data for redis `key: {member: score}`
type ZData struct {
	Key string
	Zs  []*redis.Z
}

// BatchZAdd execute `zadd` with pipeline
func (c *Cache) BatchZAdd(zDataList []ZData) error {
	pipe := c.cli.TxPipeline()
	ctx := context.TODO()

	for _, zData := range zDataList {
		key := c.genKey(zData.Key)
		// 将一个或多个成员元素及其分数值加入到有序集当中。
		// 如果某个成员已经是有序集的成员，则更新该成员的分数值，并通过重新插入该成员元素，来保证该成员在正确的位置上。
		pipe.ZAdd(ctx, key, zData.Zs...)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// ZRevRangeByScore execute `zrevrangebyscorewithscores`
func (c *Cache) ZRevRangeByScore(k string, min int64, max int64, offset int64, count int64) ([]redis.Z, error) {
	// 时间戳, 从大到小排序
	ctx := context.TODO()

	key := c.genKey(k)
	// TODO: add limit, offset, count => to ignore the too large list size
	// LIMIT 0 -1 equals no args
	// 按照分数从高到低（降序）获取有序集合中的成员，并返回这些成员的分数。
	cmds := c.cli.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:    strconv.FormatInt(min, 10),
		Max:    strconv.FormatInt(max, 10),
		Offset: offset,
		Count:  count,
	})

	return cmds.Result()
}

// BatchZRemove execute `zremrangebyscore` with pipeline
func (c *Cache) BatchZRemove(keys []string, min int64, max int64) error {
	pipe := c.cli.TxPipeline()
	ctx := context.TODO()

	minStr := strconv.FormatInt(min, 10)
	maxStr := strconv.FormatInt(max, 10)

	for _, k := range keys {
		key := c.genKey(k)
		// 移除有序集合（sorted set）中分数（score）在指定范围内的所有成员
		pipe.ZRemRangeByScore(ctx, key, minStr, maxStr)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// HashKeyField is a hash data for redis, `Key: field -> `
type HashKeyField struct {
	Key   string
	Field string
}

// Hash is a hash data  `Key: field->value`
type Hash struct {
	HashKeyField
	Value string
}

// HGet execute `hget`
func (c *Cache) HGet(hashKeyField HashKeyField) (string, error) {
	k := c.genKey(hashKeyField.Key)
	// 从哈希表（Hash）中获取指定字段（field）的值
	return c.cli.HGet(context.TODO(), k, hashKeyField.Field).Result()
}

// HSet execute `hset`
func (c *Cache) HSet(hashKeyField HashKeyField, value string) error {
	k := c.genKey(hashKeyField.Key)
	_, err := c.cli.HSet(context.TODO(), k, hashKeyField.Field, value).Result()
	return err
}

// BatchHSetWithTx execute `hset` with tx pipeline
func (c *Cache) BatchHSetWithTx(hashes []Hash) error {
	// tx, all success or all fail
	pipe := c.cli.TxPipeline()
	ctx := context.TODO()

	for _, h := range hashes {
		key := c.genKey(h.Key)
		pipe.HSet(ctx, key, h.Field, h.Value)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// BatchHGet execute `hget` with pipeline
func (c *Cache) BatchHGet(hashKeyFields []HashKeyField) (map[HashKeyField]string, error) {
	pipe := c.cli.Pipeline()

	ctx := context.TODO()
	cmds := make(map[HashKeyField]*redis.StringCmd, len(hashKeyFields))
	for _, h := range hashKeyFields {
		key := c.genKey(h.Key)
		cmd := pipe.HGet(ctx, key, h.Field)

		cmds[h] = cmd
	}

	_, err := pipe.Exec(ctx)
	// 当批量操作, 里面有个key不存在, err = redis.Nil; 但是不应该影响其他存在的key的获取
	// Nil reply returned by Redis when key does not exist.
	if err != nil && err != redis.Nil {
		return nil, err
	}

	values := make(map[HashKeyField]string, len(cmds))
	for hkf, cmd := range cmds {
		// maybe err or key missing
		// only return the HashKeyField who get value success from redis
		val, err := cmd.Result()
		if err != nil {
			continue
		} else {
			values[hkf] = val
		}
	}
	return values, nil
}

// HKeys execute `hkeys`
func (c *Cache) HKeys(hashKey string) ([]string, error) {
	key := c.genKey(hashKey)
	// 获取哈希表（Hash）中所有的字段名（field names）
	return c.cli.HKeys(context.TODO(), key).Result()
}

// Unmarshal with compress, via go-redis/cache, use s2 compression
// Note: YOU SHOULD NOT USE THE RAW msgpack.Unmarshal directly! will panic with decode fail
// 将压缩的字节切片 b 反序列化为指定的目标对象 value
func (c *Cache) Unmarshal(b []byte, value interface{}) error {
	return c.codec.Unmarshal(b, value)
}

// Marshal with compress, via go-redis/cache, use s2 compression
// Note: YOU SHOULD NOT USE THE RAW msgpack.Marshal directly!
// 将指定的对象 value 序列化为压缩的字节切片
func (c *Cache) Marshal(value interface{}) ([]byte, error) {
	return c.codec.Marshal(value)
}
