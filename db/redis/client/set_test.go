package client

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func redisSetKey(key string, val string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 设置Redis ' Set key value [expiration] '命令。
	err := rdb.Set(ctx, key, val, time.Hour).Err()
	if err != nil {
		fmt.Printf("redis set failed, err: %v\n", err)
		return err
	}
	return nil
}

func TestSet(t *testing.T) {
	if err := initRedisV9Client(); err != nil {
		fmt.Printf("initRedisV9Client failed: %v\n", err)
		return
	}
	fmt.Println("initRedisV9Client started successfully")
	defer rdb.Close()

	// set value
	err := redisSetKey("name", "xia")
	if err != nil {
		fmt.Printf("redisSetKey failed: %v\n", err)
		return
	}
	fmt.Println("redisSetKey succeeded")
}
