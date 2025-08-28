package main

import (
	"testing"

	"trpc.group/trpc-go/trpc-database/goredis"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
)

// TestNewUniversalClient 原生client用例
func TestNewUniversalClient(t *testing.T) {
	// 扩展接口 (use redis.NewUniversalClient)
	target := "redis://127.0.0.1:6379/0"
	cli, err := goredis.New("trpc.gamecenter.test.redis", client.WithTarget(target))
	// TODO 当前版本不支持 target 参数，需要手动设置
	// cli, err := goredis.New("trpc.gamecenter.test.redis")
	if err != nil {
		t.Fatalf("new fail err=[%v]", err)
	}

	// go-redis Set
	result, err := cli.Set(trpc.BackgroundContext(), "key", "value", 0).Result()
	if err != nil {
		t.Fatalf("set fail err=[%v]", err)
	}
	t.Logf("Set result=[%v]", result)

	// go-redis Get
	value, err := cli.Get(trpc.BackgroundContext(), "key").Result()
	if err != nil {
		t.Fatalf("get fail err=[%v]", err)
	}
	t.Logf("Get value=[%v]", value)
}
