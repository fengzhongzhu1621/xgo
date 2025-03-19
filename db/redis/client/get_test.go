package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	redisV9 "github.com/redis/go-redis/v9"
)

func redisGetKey(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redisV9.Nil {
			return "", nil
			// DeadlineExceeded是Context返回的错误。当上下文的截止日期过去时发生错误。
		} else if err == context.DeadlineExceeded {
			return "", fmt.Errorf("获取值超时")
		} else {
			return "", fmt.Errorf("获取值失败: %v", err)
		}
	}

	if val == "" {
		return "", nil
	}

	return val, nil
}

func TestGet(t *testing.T) {
	if err := initRedisV9Client(); err != nil {
		fmt.Printf("initRedisV9Client failed: %v\n", err)
		return
	}
	fmt.Println("initRedisV9Client started successfully")
	defer rdb.Close() // Close 关闭客户端，释放所有打开的资源。关闭客户端是很少见的，因为客户端是长期存在的，并在许多例程之间共享。

	// TODO redisCommand()

	// get key
	value, _ := redisGetKey("key")
	fmt.Printf("get key: %v\n", value)
}
