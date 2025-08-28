package redislock

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type IRedisClient interface {
	AcquireLock(ctx context.Context, key string, ttl time.Duration, opt *IOptions) (*Lock, error)
	// SetNX 置如果不存在
	SetNX(
		ctx context.Context,
		key string,
		value interface{},
		expiration time.Duration,
	) *redis.BoolCmd
	// Eval 执行Lua脚本
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	// EvalSha 执行Lua脚本
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
	// ScriptExists 检查脚本是否存在
	ScriptExists(ctx context.Context, scripts ...string) *redis.BoolSliceCmd
	// ScriptLoad 加载脚本到Redis
	ScriptLoad(ctx context.Context, script string) *redis.StringCmd
	// Get 获取键的值
	Get(ctx context.Context, key string) *redis.StringCmd
	// Del 删除键
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	// CheckUnLock 检查一个特定的键（key）是否已经被解锁
	CheckUnLock(ctx context.Context, key string) (bool, error)
	// ForceRelease 强制释放一个特定的键（key），无论该键是否已经被当前客户端锁定
	ForceRelease(ctx context.Context, key, value string) (bool, error)
}

type RedislockClient struct {
	client IRedisClient
	tmp    []byte
	tmpMu  sync.Mutex
}

type Lock struct {
	client *RedislockClient
	key    string
	value  string
}

func (l *Lock) Key() string {
	return l.key
}

func (l *Lock) Token() string {
	return l.value[:22]
}

func (l *Lock) Metadata() string {
	return l.value[22:]
}

// TTL 获取 key 的 剩余生存时间（TTL）
func (l *Lock) TTL(ctx context.Context) (time.Duration, error) {
	// pttl 获取 key 的 剩余生存时间（TTL）
	res, err := luaPTTL.Run(ctx, l.client.client, []string{l.key}, l.value).Result()
	if err == redis.Nil {
		// 键不存在于 Redis 中
		return 0, nil
	} else if err != nil {
		// 获取 TTL 时发生了其他错误
		return 0, err
	}

	// 键存在且有一个剩余生存时间，获得正确的 TTL 单位
	if num := res.(int64); num > 0 {
		return time.Duration(num) * time.Second, nil
	}

	// -3 键的值不匹配 or TTL 为 0
	return 0, nil
}

// Lease 刷新 key 的过期时间（TTL）
func (l *Lock) Lease(ctx context.Context, ttl time.Duration, opt *IOptions) error {
	// 刷新 Redis 中特定键（KEYS[1]）的过期时间（TTL）
	ttlVal := strconv.FormatInt(int64(ttl/time.Second), 10)
	status, err := luaRefresh.Run(ctx, l.client.client, []string{l.key}, l.value, ttlVal).Result()
	if err != nil {
		// 刷新 TTL 时发生了其他错误
		return err
	} else if status == int64(1) {
		// 键的 TTL 成功刷新
		return nil
	}
	return ErrNotObtained
}

// Release 释放 Redis 中特定键（KEYS[1]）的锁
func (l *Lock) Release(ctx context.Context) error {
	res, err := luaRelease.Run(ctx, l.client.client, []string{l.key}, l.value).Result()
	if err == redis.Nil {
		// 键不存在于 Redis 中，锁没有被当前客户端持有
		return ErrLockNotHeld
	} else if err != nil {
		// 在释放锁时发生了其他错误
		return err
	}

	// 锁存在，但是没有被当前客户端持有
	if i, ok := res.(int64); !ok || i != 1 {
		return ErrLockNotHeld
	}

	// 锁成功释放
	return nil
}

func NewRedislockClient(client IRedisClient) *RedislockClient {
	return &RedislockClient{
		client: client,
	}
}

// randomToken 生成一个随机的令牌，这个令牌将用作分布式锁的值
func (c *RedislockClient) randomToken() (string, error) {
	c.tmpMu.Lock()
	defer c.tmpMu.Unlock()

	// 创建一个新的16字节切片
	if len(c.tmp) == 0 {
		c.tmp = make([]byte, 16)
	}

	// 生成随机数据填充tmp切片
	if _, err := io.ReadFull(rand.Reader, c.tmp); err != nil {
		return "", err
	}

	// 对生成的随机数据进行编码，得到一个URL安全的字符串作为令牌返回。
	// URL safe 编码, 替换掉字符串中的特殊字符，+ 和 /，末尾不补 =
	return base64.RawURLEncoding.EncodeToString(c.tmp), nil
}

// AcquireLock 加锁（支持重试）
func (c *RedislockClient) AcquireLock(
	ctx context.Context,
	key string,
	ttl time.Duration,
	opt *IOptions,
) (*Lock, error) {
	token, err := c.randomToken()
	if err != nil {
		return nil, err
	}

	// 用于获取 Options 结构体中的 Metadata 字段的值
	value := token + opt.GetMetadata()
	// 获取下一次重试的等待时间策略
	retry := opt.GetRetryStrategy()

	// 该上下文会在指定的截止时间到达时自动取消
	deadlinectx, cancel := context.WithDeadline(ctx, time.Now().Add(ttl))
	defer cancel()

	var timer *time.Timer
	for {
		// 设置 key 的过期时间（分布式环境唯一性）
		ok, err := c.acquire(deadlinectx, key, value, ttl)
		if err != nil {
			return nil, err
		} else if ok {
			// 成功设计
			return &Lock{client: c, key: key, value: value}, nil
		}

		// key 已经存在，已被其它进程加锁，获得下一次重新加锁的间隔
		backoff := retry.NextBackoff()
		if backoff < 1 {
			// 不会重新加锁
			return nil, ErrNotObtained
		}

		// 在指定间隔后重新触发
		if timer == nil {
			timer = time.NewTimer(backoff)
			defer timer.Stop()
		} else {
			// 重置定时器的到期时间，定时器会重新开始计时，并在新的到期时间到达时触发
			timer.Reset(backoff)
		}

		// 阻塞等待
		select {
		case <-deadlinectx.Done():
			return nil, ErrNotObtained
		case <-timer.C:
		}
	}
}

// acquire 设置键值对，但只有在键不存在时才会执行
func (c *RedislockClient) acquire(
	ctx context.Context,
	key, value string,
	ttl time.Duration,
) (bool, error) {
	return c.client.SetNX(ctx, key, value, ttl).Result()
}

// CheckUnLock 检查一个特定的键（key）是否已经被解锁
func (c *RedislockClient) CheckUnLock(ctx context.Context, key string) (bool, error) {
	// 根据 key 获取加锁信息
	_, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// 键不存在于 Redis 中。在这种情况下，我们认为键已经解锁，因此返回 true 和 nil 错误
		return true, nil
	}
	if err != nil {
		// 获取键值时发生了其他错误。在这种情况下，我们返回 false 和该错误对象。
		return false, err
	}

	// 键存在于 Redis 中。在这种情况下，我们认为键仍然被锁定
	return false, nil
}

// ForceRelease 强制释放一个特定的键（key），无论该键是否已经被当前客户端锁定
func (c *RedislockClient) ForceRelease(ctx context.Context, key, value string) (bool, error) {
	// 尝试从 Redis 中删除指定的键
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		// 删除键时发生了错误
		return false, err
	}

	// 键成功删除
	return true, nil
}
