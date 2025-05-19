package ratelimiter

import (
	"container/list" // 双向链表
	"sync"
	"time"
)

// SlidingWindowCounter 实现了一个滑动窗口计数器(Sliding Window Counter)算法，用于限流(Rate Limiting)。
// 与固定窗口计数器不同，滑动窗口计数器能更精确地控制单位时间内的请求数量。
// * 滑动窗口：时间窗口是滑动的，随着新请求的到来而移动
// * 精确限流：比固定窗口更精确地控制单位时间内的请求数量
// * 并发安全：通过互斥锁保证多线程环境下的正确性
// * 内存消耗：需要存储每个请求的时间戳，可能消耗较多内存(在高并发场景下)
type SlidingWindowCounter struct {
	mu     sync.Mutex    // 互斥锁，保证并发安全
	events *list.List    // 双向链表，存储请求时间戳，链表长度等于当前窗口内的请求数，链表的值是请求时间
	limit  int           // 窗口内允许的最大请求数
	window time.Duration // 滑动时间窗口长度
}

func NewSlidingWindowCounter(limit int, window time.Duration) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		events: list.New(), // 初始化空链表
		limit:  limit,
		window: window,
	}
}

func (sw *SlidingWindowCounter) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()

	// 清理过期的事件，去掉滑动窗口左侧的历史请求
	for sw.events.Len() > 0 {
		// 检查链表头部的时间戳是否已过期（滑动窗口过期）（左侧）
		if sw.events.Front().Value.(time.Time).Add(sw.window).Before(now) {
			// 移除过期时间戳（从链表中删除）
			sw.events.Remove(sw.events.Front())
		} else {
			// 遇到未过期的时间戳，停止清理
			break
		}
	}

	// 检查当前窗口内的请求数是否超过限制，在滑动窗口右侧添加新的请求
	if sw.events.Len() < sw.limit {
		sw.events.PushBack(now) // 添加当前时间戳到链表尾部（右侧）
		return true             // 允许请求
	}

	// 拒绝请求
	return false
}
