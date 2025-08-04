package stream

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestXAdd(t *testing.T) {
	// Redis 客户端初始化
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	// 使用带超时的上下文
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 事件数据准备
	event := map[string]interface{}{"message": "Critical alert! Server down."}

	// 发布消息到 Redis Stream
	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: "alerts", // 目标 Stream 名称 ("alerts")
		Values: event,    // 要发布的消息内容（即之前准备的 event map）
		MaxLen: 1000,     // 限制 Stream 最多 1000 条消息，防止 Stream 无限增长
	}).Result()
	if err != nil {
		log.Fatalf("发布事件失败: %v", err)
	}
	fmt.Println("事件发布成功")
}
