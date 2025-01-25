package redsync

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func TestLockWithExpiry(t *testing.T) {
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

	// 模拟工作
	fmt.Println("Lock acquired! Performing work...")
	time.Sleep(5 * time.Second)

	// 释放锁。一般来说，锁会在持有锁的操作完成后被释放。
	// 使用 defer 确保锁被释放
	defer func() {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			log.Printf("无法释放锁: ok=%v, err=%v", ok, err)
		} else {
			log.Println("锁已成功释放")
		}
	}()
}
