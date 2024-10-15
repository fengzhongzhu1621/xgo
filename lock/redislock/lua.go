package redislock

import (
	"errors"

	"github.com/go-redis/redis/v8"
)

var (
	// PEXPIRE 和 EXPIRE 是 Redis 中用于设置键的过期时间的两个命令，它们的主要区别在于时间单位的差异
	// pexpire 的 lua 实现，刷新 Redis 中特定键（KEYS[1]）的过期时间（TTL）
	// 使用 redis.call("get", KEYS[1]) 获取 Redis 中指定键（KEYS[1]）的值。
	// 检查获取到的值是否等于传入的参数（ARGV[1]）。
	// 如果值相等，则使用 redis.call("expire", KEYS[1], ARGV[2]) 刷新键的过期时间（TTL），并将其返回。
	// 如果值不相等，则返回 0，表示键的值不匹配。
	luaRefresh = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("expire", KEYS[1], ARGV[2]) else return 0 end`)

	// pttl 的 lua 实现
	// 使用 redis.call("get", KEYS[1]) 获取 Redis 中指定键（KEYS[1]）的值。
	// 检查获取到的值是否等于传入的参数（ARGV[1]）。
	// 如果值相等，则使用 redis.call("ttl", KEYS[1]) 获取键的剩余生存时间（TTL），并将其返回。
	// 如果值不相等，则返回 -3，表示键的值不匹配。
	luaPTTL = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("ttl", KEYS[1]) else return -3 end`)

	// 用于释放 Redis 中特定键（KEYS[1]）的锁
	// 使用 redis.call("get", KEYS[1]) 获取 Redis 中指定键（KEYS[1]）的值。
	// 检查获取到的值是否等于传入的参数（ARGV[1]）。
	// 如果值相等，则使用 redis.call("del", KEYS[1]) 删除键，并返回删除操作的结果。
	// 如果值不相等，则返回 0，表示键的值不匹配。
	luaRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)

	// 加锁冲突，且不会重复加锁
	ErrNotObtained = errors.New("redislock: not obtained")

	// 锁没有被当前客户端持有
	ErrLockNotHeld = errors.New("redislock: lock not held")
)
