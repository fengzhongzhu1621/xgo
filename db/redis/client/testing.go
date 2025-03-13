package client

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

// NewTestRedisClient 创建一个新的 Redis 客户端实例，该实例连接到由 miniredis 模拟的 Redis 服务器
// returns: 一个指向 redis.Client 结构体的指针，表示 Redis 客户端实例
func NewTestRedisClient() *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}
