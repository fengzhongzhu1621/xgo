package redsync

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var redisClient *redis.Client

func initRedisClient(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr, // Redis 服务器地址
	})

	// 测试连接
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("无法连接到 Redis: %v", err)
	}

	return client
}

func TestLock(t *testing.T) {
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

	// 获取锁
	mutex := rs.NewMutex("mylock")

	// 获取锁。如果锁已经被其他实例持有，当前实例会阻塞，直到获取到锁。
	// 会一直阻塞，直到成功获取锁。
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
