package client

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func hGetAll(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	val1 := rdb.HMGet(ctx, "user", "name", "age").Val()
	fmt.Printf("redis HMGet %v\n", val1)

	val2 := rdb.HGet(ctx, "user", "age").Val()
	fmt.Printf("redis HGet value: %v\n", val2)

	val3, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		// redis.Nil
		// 其它错误
		fmt.Printf("hgetall failed, err: %v\n", err)
		return nil, err
	}
	return val3, nil
}

func TestHGetAll(t *testing.T) {
	if err := initRedisV9Client(); err != nil {
		fmt.Printf("initRedisV9Client failed: %v\n", err)
		return
	}
	fmt.Println("initRedisV9Client started successfully")
	defer rdb.Close() // Close 关闭客户端，释放所有打开的资源。关闭客户端是很少见的，因为客户端是长期存在的，并在许多例程之间共享。

	// hgetall()
	value, err := hGetAll("user")
	if err != nil {
		fmt.Printf("hGetAll failed with error: %v\n", err)
		return
	}
	fmt.Printf("hGetAll successful, value: %v\n", value)
}
