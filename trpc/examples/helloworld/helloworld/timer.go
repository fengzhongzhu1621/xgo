package main

import (
	"context"
	"os"
	"time"

	"trpc.group/trpc-go/trpc-go/log"
)

func handleLocalTimer(ctx context.Context) error {
	log.Info("do local timer processing...")
	// 日志支持 log.Field
	log.With(log.Field{Key: "user", Value: os.Getenv("USER")}).Info("test log.Field ok")
	// 显式地获取默认 Logger
	// l := log.Get("default")
	l := log.GetDefaultLogger()
	n := 1
	fields := make([]log.Field, 0, n)
	for i := 0; i < n; i++ {
		fields = append(fields, log.Field{Key: "key", Value: os.Getenv("USER")})
	}
	l = l.With(fields...)
	l.Info("test GetDefaultLogger() ok")
	// 使用 context 级别的日志打印
	log.InfoContext(ctx, "do local timer processing end")

	return nil
}

func handleDistributedTimer(_ context.Context) error {
	log.Info("do distributed timer processing...")
	return nil
}

// Scheduler 调度策略
type scheduler struct {
}

// Schedule 自己通过数据存储实现互斥任务定时器
func (s *scheduler) Schedule(serviceName string, newNode string, holdTime time.Duration) (nowNode string, err error) {
	// 使用 Redis 的 setnx 来实现
	// setnx serviceName newNode ex expireSeconds
	// if fail, nowNode = get serviceName
	return nowNode, nil
}
