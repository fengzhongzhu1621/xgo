package ratelimiter

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestAllow(t *testing.T) {
	// 创建一个每秒产生1个令牌，最多存储5个令牌的限制器
	r := rate.NewLimiter(1, 5)

	// 模拟10个事件
	for i := 0; i < 10; i++ {
		if r.Allow() {
			fmt.Println("Handle event", i)
		} else {
			fmt.Println("Rate limited", i)
		}
		// 为了更好地观察速率限制效果，可以在这里添加短暂的延迟
		// time.Sleep(200 * time.Millisecond)
	}
}

func TestAllowN(t *testing.T) {
	r := rate.NewLimiter(1, 5) // 每秒放置1个令牌，最多存储5个令牌
	for i := 0; i < 10; i++ {
		if r.AllowN(time.Now(), 2) { // 每次处理2个事件
			fmt.Println("Handle events", i, i+1)
		} else {
			fmt.Println("Rate limited", i, i+1)
		}
		time.Sleep(time.Second)
	}
}

func TestReserve(t *testing.T) {
	// 创建一个每秒产生1个令牌，最多存储5个令牌的限制器
	r := rate.NewLimiter(1, 5)

	// 模拟10个事件
	for i := 0; i < 10; i++ {
		res := r.Reserve() // 预订一个时间段
		if res.OK() {      // 检查预订是否成功
			fmt.Println("Handle event", i)
			time.Sleep(res.Delay()) // 等待预订的时间段
		} else {
			fmt.Println("Rate limited", i)
		}
	}
}
