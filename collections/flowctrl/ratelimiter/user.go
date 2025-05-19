package ratelimiter

import (
	"sync"
	"time"
)

// UserRateLimiter 用户速率限制器
type UserRateLimiter struct {
	mu         sync.Mutex           // 互斥锁，保证并发安全
	userLimits map[string]int       // 每个用户的速率限制阈值
	userCounts map[string]int       // 每个用户当前的请求计数
	limit      int                  // 全局限制阈值（如果需要全局限制）
	window     time.Duration        // 时间窗口，单位为纳秒
	resetTimes map[string]time.Time // 每个用户的上次重置时间
}

// NewUserRateLimiter 创建一个新的用户速率限制器实例
//
// 参数:
//   - limit: 每个用户在时间窗口内允许的最大请求次数
//   - window: 时间窗口的持续时间，例如 time.Minute 表示每分钟
//
// 返回值:
//   - *UserRateLimiter: 初始化后的用户速率限制器实例
func NewUserRateLimiter(limit int, window time.Duration) *UserRateLimiter {
	return &UserRateLimiter{
		userLimits: make(map[string]int),
		userCounts: make(map[string]int),
		limit:      limit,
		window:     window,
		resetTimes: make(map[string]time.Time),
	}
}

// Allow 尝试允许指定用户的请求
//
// 参数:
//   - userId: 用户的唯一标识符
//
// 返回值:
//   - bool: 如果请求被允许返回 true，否则返回 false
func (url *UserRateLimiter) Allow(userId string) bool {
	url.mu.Lock()
	defer url.mu.Unlock()

	// 初始化用户的重置时间和计数（如果尚未初始化）
	now := time.Now()
	if _, exists := url.resetTimes[userId]; !exists {
		url.resetTimes[userId] = now
	}

	// 检查是否需要重置计数
	if now.Sub(url.resetTimes[userId]) >= url.window {
		url.userCounts[userId] = 0
		url.resetTimes[userId] = now
	}

	// 检查是否超过限制
	if url.userCounts[userId] < url.limit {
		url.userCounts[userId]++
		return true
	}

	return false
}

// SetLimit 设置指定用户的速率限制阈值
//
// 参数:
//   - userId: 用户的唯一标识符
//   - newLimit: 新的速率限制阈值
func (url *UserRateLimiter) SetLimit(userId string, newLimit int) {
	url.mu.Lock()
	defer url.mu.Unlock()

	url.userLimits[userId] = newLimit
}

// GetCount 获取指定用户的当前请求计数
//
// 参数:
//   - userId: 用户的唯一标识符
//
// 返回值:
//   - int: 当前请求计数
func (url *UserRateLimiter) GetCount(userId string) int {
	url.mu.Lock()
	defer url.mu.Unlock()

	return url.userCounts[userId]
}

// Reset 重置指定用户的请求计数和重置时间（主要用于测试）
//
// 参数:
//   - userId: 用户的唯一标识符
func (url *UserRateLimiter) Reset(userId string) {
	url.mu.Lock()
	defer url.mu.Unlock()

	url.userCounts[userId] = 0
	url.resetTimes[userId] = time.Now()
}

// ResetAll 重置所有用户的请求计数和重置时间（主要用于测试）
func (url *UserRateLimiter) ResetAll() {
	url.mu.Lock()
	defer url.mu.Unlock()

	now := time.Now()
	for userId := range url.resetTimes {
		url.userCounts[userId] = 0
		url.resetTimes[userId] = now
	}
}
