package redsync

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func TestTryLock(t *testing.T) {
	// 创建 Redis 客户端
	redisClient = initRedisClient("localhost:6379")

	// 确保在程序结束时关闭 Redis 客户端
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("关闭 Redis 客户端时出错: %v", err)
		}
	}()

	// 创建 Redsync 实例
	pool := goredis.NewPool(redisClient)
	rs := redsync.New(pool)

	// 获取锁, 如果锁在 5 秒内没有释放，其他实例将可以获取该锁。
	mutex := rs.NewMutex("mylock", redsync.WithExpiry(5*time.Second))

	// 获取锁。如果锁已经被其他实例持有，当前实例会阻塞，直到获取到锁。
	if err := mutex.Lock(); err != nil {
		log.Fatalf("无法获取锁: %v", err)
	}

	// 尝试获取锁，如果无法获取则立即返回
	if err := mutex.TryLock(); err != nil {
		if err != nil {
			log.Fatalf("could not acquire lock: %v", err)
		}
		log.Println("Lock not acquired")
	} else {
		// 成功获取锁
		fmt.Println("Lock acquired!")
		time.Sleep(5 * time.Second) // 模拟工作

		// 释放锁
		if ok, err := mutex.Unlock(); !ok || err != nil {
			log.Printf("无法释放锁: ok=%v, err=%v", ok, err)
		} else {
			log.Println("锁已成功释放")
		}

		fmt.Println("Lock released!")
	}
}
