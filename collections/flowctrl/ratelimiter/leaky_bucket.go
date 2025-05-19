package ratelimiter

import (
	"sync"
	"time"
)

// LeakyBucket 实现了一个漏桶(Leaky Bucket)限流算法，用于控制请求的处理速率。
// * 平滑限流：漏桶算法会平滑处理突发流量，即使短时间内有大量请求到来，也会按照固定速率处理
// * 固定速率：无论请求量如何变化，处理速率始终固定为1/rate个请求/秒
// * 并发安全：通过互斥锁保证多线程环境下的正确性
// * 内存效率：只需要存储少量状态信息
type LeakyBucket struct {
	mu        sync.Mutex    // 互斥锁，保证并发安全
	capacity  int           // 桶的容量（最大令牌数）
	available int           // 当前可用令牌数
	rate      time.Duration // 令牌泄漏速率（每个令牌需要多长时间） QPS
	lastTime  time.Time     // 上次泄漏令牌的时间
}

func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity:  capacity,   // 设置桶的最大容量
		available: capacity,   // 初始时桶是满的
		rate:      rate,       // 设置令牌泄漏速率（固定速率的漏桶），多长时间处理一个请求，即请求处理的QPS
		lastTime:  time.Now(), // 记录初始时间
	}
}

func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 计算自上次操作以来的时间（计算漏水的时间间隔）
	now := time.Now()
	elapsed := now.Sub(lb.lastTime)

	// 追加可用令牌
	lb.available += int(elapsed / lb.rate) // 根据时间流逝增加令牌
	if lb.available > lb.capacity {        // 确保不超过桶容量
		lb.available = lb.capacity
	}
	lb.lastTime = now

	// 检查是否有可用令牌
	if lb.available > 0 {
		lb.available-- // 消耗一个令牌
		return true    // 允许请求
	}

	// 拒绝请求
	return false
}
