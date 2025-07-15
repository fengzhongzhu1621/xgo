package main

import (
	"context"
	"time"

	"github.com/silenceper/log"
)

func handleLocalTimer(_ context.Context) error {
	log.Info("do local timer processing...")
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
