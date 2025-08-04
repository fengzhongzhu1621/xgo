package ratelimiter

import (
	"testing"
	"time"

	"github.com/RussellLuo/slidingwindow"
)

// 一个使用滑动窗口算法（Sliding Window）实现的限流器（Rate Limiter）测试用例
// 通过维护两个窗口（当前窗口和前一个窗口）的请求计数，估算当前时间点的请求速率，从而实现滑动窗口限流。
func TestSlidingWindow(t *testing.T) {
	// 创建了一个本地滑动窗口实例
	// Window  接口定义了窗口的基本操作，包括添加计数、获取计数、重置窗口等。
	// 提供了两种实现：
	// LocalWindow：本地窗口，适用于单机应用。
	// SyncWindow：同步窗口，适用于分布式环境。
	window, stop := slidingwindow.NewLocalWindow()
	defer stop()

	// 创建一个限制器，限制每分钟最多 100 次请求
	limiter, _ := slidingwindow.NewLimiter(
		time.Minute,
		100,
		func() (slidingwindow.Window, slidingwindow.StopFunc) {
			return window, stop
		},
	)

	// 测试前100次请求应该被允许
	for i := 0; i < 100; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d should be allowed, but was denied", i)
		}
	}

	// 测试第101次请求应该被拒绝
	if limiter.Allow() {
		t.Error("Request 101 should be denied, but was allowed")
	}
}
