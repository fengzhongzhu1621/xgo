package stream

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestXRead(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	ctx := context.Background()

	for {
		res, err := client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"alerts", "$"}, // "$" 表示从最新位置开始
			Count:   1,
			Block:   0, // 非阻塞等待
		}).Result()
		if err != nil {
			log.Fatalf("读取事件失败: %v", err)
		}

		for _, stream := range res {
			for _, msg := range stream.Messages {
				fmt.Printf("处理事件: %v\n", msg.Values)
			}
		}
	}
}

func TestXReadGroup(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	// 使用带超时的上下文
	ctx := context.Background()

	// 创建消费者组（如果不存在）
	_, err := client.XGroupCreateMkStream(ctx, "alerts", "my-group", "$").Result()
	if err != nil && !errors.Is(err, redis.TxFailedErr) { // 忽略 "BUSYGROUP" 错误（组已存在）
		log.Fatalf("创建消费者组失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 设置优雅退出
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		cancel()
	}()

	// 3. 消费消息（使用 XReadGroup）
	for {
		select {
		case <-ctx.Done():
			fmt.Println("程序退出")
			return
		default:
			// "$"：表示从 Stream 的最新位置（最新消息）开始读取
			res, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    "my-group",
				Consumer: "consumer-1",
				Streams:  []string{"alerts", ">"}, // ">" 表示从未消费的消息开始
				Count:    1,
				Block:    1000, // 阻塞 1 秒等待新消息
			}).Result()
			if err != nil {
				log.Printf("读取事件失败: %v", err)
				continue // 继续尝试
			}

			for _, stream := range res {
				for _, msg := range stream.Messages {
					fmt.Printf("处理事件: %v\n", msg.Values)

					// 4. 确认消息（防止重复消费）
					confirmedCount, err := client.XAck(ctx, "alerts", "my-group", msg.ID).Result()
					if err != nil {
						log.Printf("确认消息失败: %v", err)
						// 可能需要重试或记录错误
					} else {
						log.Printf("成功确认 %d 条消息", confirmedCount)
						// 通常 confirmedCount == 1
					}

				}
			}
		}
	}
}
