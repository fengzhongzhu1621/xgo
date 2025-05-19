package ratelimiter

import (
	"sync"
	"time"
)

// FixedWindowCounter 一个固定窗口计数器(Fixed Window Counter)算法，用于限流(Rate Limiting)。
// * 固定窗口：时间窗口是固定的，不滑动
// * 简单高效：实现简单，计算量小
// * 并发安全：通过互斥锁保证多线程环境下的正确性
// * 边界问题：在窗口切换瞬间可能出现突发流量(例如窗口末尾和开始瞬间允许的请求数可能超过限制)
type FixedWindowCounter struct {
	mu           sync.Mutex    // 互斥锁，保证并发安全
	requestCount int           // 当前窗口内的请求计数
	limit        int           // 窗口内允许的最大请求数
	window       time.Duration // 时间窗口长度
	resetTime    time.Time     // 当前窗口重置的时间点
}

// NewFixedWindowCounter 创建一个固定窗口计数器限流器
// counter := NewFixedWindowCounter(100, time.Second) // 每秒最多100个请求
func NewFixedWindowCounter(limit int, window time.Duration) *FixedWindowCounter {
	return &FixedWindowCounter{
		limit:     limit,  // 窗口内允许的最大请求数
		window:    window, // 时间窗口长度
		resetTime: time.Now(),
	}
}

func (fw *FixedWindowCounter) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()

	// 检查是否需要重置计数器
	if now.Sub(fw.resetTime) >= fw.window {
		fw.requestCount = 0 // 重置计数器
		fw.resetTime = now  // 更新重置时间
	}

	// 检查是否超过限制
	if fw.requestCount < fw.limit {
		fw.requestCount++ // 增加计数
		return true       // 允许请求
	}

	// 拒绝请求
	return false
}
