package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestPipeline(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 创建管道
	pipe := rdb.Pipeline()

	// 添加命令到管道
	inc := pipe.Incr(ctx, "counter")        // 自增 counter
	set := pipe.Set(ctx, "key", "value", 0) // 设置 key=value
	get := pipe.Get(ctx, "key")             // 获取 key 的值

	// 执行管道
	_, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}

	// 查看结果
	fmt.Println("counter 自增后:", inc.Val()) // counter 自增后: 1
	fmt.Println("set 命令错误:", set.Err())    // set 命令错误: <nil>
	fmt.Println("key 的值:", get.Val())      // key 的值: value
}
