package ratelimiter

import (
	"sync"
	"time"
)

// TokenBucket 实现了一个令牌桶(Token Bucket)限流算法，用于控制请求的处理速率。
type TokenBucket struct {
	mu       sync.Mutex    // 互斥锁，保证并发安全
	capacity int           // 桶的容量，即最大令牌数
	tokens   int           // 当前令牌数量
	rate     time.Duration // 令牌生成速率，每 rate 时间生成一个令牌
	lastTime time.Time     // 上次填充令牌的时间
}

func NewTokenBucket(capacity int, rate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		tokens:   capacity,
		rate:     rate,
		lastTime: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 计算时间间隔
	now := time.Now()
	elapsed := now.Sub(tb.lastTime)

	// 计算时间间隔内生成的可用令牌数
	tb.tokens += int(elapsed / tb.rate)
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastTime = now

	// 申请一个令牌
	// 流出速度不确定，一旦有突增的流量，令牌桶里已有的令牌可以短暂的应对突发流量
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// AllowN 尝试从令牌桶中获取 n 个令牌
//
// 参数:
//   - n: 需要获取的令牌数量
//
// 返回值:
//   - bool: 如果获取成功返回 true，否则返回 false
func (tb *TokenBucket) AllowN(n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastTime)

	// 计算在 elapsed 时间内可以生成的令牌数量
	newTokens := int(elapsed / tb.rate)

	// 更新令牌数量，但不能超过容量
	tb.tokens += newTokens
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	// 更新上次填充时间
	tb.lastTime = now

	// 检查是否有足够的令牌
	// 申请 n 个令牌
	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}

	return false
}

// SetRate 动态调整令牌生成速率
func (tb *TokenBucket) SetRate(rate time.Duration) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.rate = rate
}

// SetCapacity 动态调整桶的容量
func (tb *TokenBucket) SetCapacity(capacity int) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.capacity = capacity
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
}
