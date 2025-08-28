package ratelimiter

import (
	"fmt"
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

// LeakyBucket2 结构体，包含请求队列
type LeakyBucket2 struct {
	queue chan struct{} // 请求队列
}

func NewLeakyBucket2(capacity int) *LeakyBucket2 {
	return &LeakyBucket2{
		queue: make(chan struct{}, capacity),
	}
}

// push 将请求放入队列，如果队列满了，返回 false，表示请求被丢弃
func (lb *LeakyBucket2) push() bool {
	// 如果通道可以发送，请求被接受
	select {
	case lb.queue <- struct{}{}:
		return true
	default:
		return false
	}
}

// process 从队列中取出请求并模拟处理过程
// 启动请求处理循环    go lb.process()
func (lb *LeakyBucket2) process() {
	for range lb.queue {
		// 使用 range 来持续接收队列中的请求，以恒定速率从桶中取出请求进行处理
		fmt.Println("Request processed at", time.Now().Format("2006-01-02 15:04:05"))
		// 模拟请求处理时间
		time.Sleep(100 * time.Millisecond)
	}
}
