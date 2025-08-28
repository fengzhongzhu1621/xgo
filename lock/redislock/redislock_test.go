package redislock

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

func TestRedislock(t *testing.T) {
	// 创建一个 Redis 客户端，连接到本地 Redis 服务（默认端口 6379）
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	})
	defer client.Close()

	// 创建一个基于 Redis 的分布式锁管理器
	locker := redislock.New(client)

	// 尝试获取名为 "my-key" 的锁，锁的过期时间为 100ms
	ctx := context.Background()
	lock, err := locker.Obtain(ctx, "my-key", 100*time.Millisecond, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		// 如果锁已被占用（ErrNotObtained），打印提示信息
		fmt.Println("Could not obtain lock!")
	} else if err != nil {
		log.Fatalln(err)
	}
	// 如果获取成功，defer lock.Release(ctx) 确保测试结束后释放锁
	defer lock.Release(ctx)

	// 检查锁是否仍然有效

	fmt.Println("I have a lock!")
	// 等待 50ms（锁剩余 50ms 有效期）。
	time.Sleep(50 * time.Millisecond)
	// 检查锁的剩余生存时间（TTL）
	if ttl, err := lock.TTL(ctx); err != nil {
		log.Fatalln(err)
	} else if ttl > 0 {
		// 如果 TTL > 0，说明锁仍然有效
		fmt.Println("Yay, I still have my lock!")
	}

	// 续约锁：调用 lock.Refresh() 延长锁的有效期（再延长 100ms），注意不是追加，而是重置
	// 当调用 lock.Refresh(ctx, newTTL, options) 时，Redis 会更新该锁的过期时间（TTL）为 newTTL。
	// Refresh 是原子操作，Redis 会确保在更新 TTL 时不会发生竞争条件（即不会出现两个客户端同时修改锁的情况）。
	if err := lock.Refresh(ctx, 100*time.Millisecond, nil); err != nil {
		log.Fatalln(err)
	}

	// 新的锁的过期时间是100ms，等待100ms后过期
	time.Sleep(100 * time.Millisecond)
	if ttl, err := lock.TTL(ctx); err != nil {
		log.Fatalln(err)
	} else if ttl == 0 {
		fmt.Println("Now, my lock has expired!")
	}
}

func TestRedisLockWithRetry(t *testing.T) {
	// 连接 Redis
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	defer client.Close()

	// 创建锁客户端
	locker := redislock.New(client)

	// 配置重试策略（线性退避，每次等待 500ms）
	opts := &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(500 * time.Millisecond),
	}

	ctx := context.Background()
	// 尝试获取锁（最多重试直到成功或超时）
	// 如果锁被其他客户端持有，当前客户端不会立即失败，而是按照 500ms 的间隔重试。
	lock, err := locker.Obtain(ctx, "my-key", 10*time.Second, opts)
	if err == redislock.ErrNotObtained {
		// 明确表示未获取到锁（可能因重试次数耗尽或超时）。
		fmt.Println("无法获取锁！")
	} else if err != nil {
		// 其他错误（如 Redis 连接问题）
		log.Fatalf("获取锁失败: %v", err)
	}
	defer lock.Release(ctx) // 确保释放锁

	fmt.Println("成功获取锁！")

	// 模拟业务逻辑
	time.Sleep(2 * time.Second)
}
