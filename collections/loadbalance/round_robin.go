package loadbalance

import (
	"errors"
	"sync"
)

// RoundRobinBalance 轮询负载均衡器
type RoundRobinBalance struct {
	curIndex int      // 当前索引
	rss      []string // 服务器地址列表
	mu       sync.Mutex
}

// Add 添加服务器地址
// 参数：可变参数，至少需要一个服务器地址
// 返回值：如果参数不足，返回错误
func (r *RoundRobinBalance) Add(addrs ...string) error {
	if len(addrs) == 0 {
		return errors.New("at least 1 parameter is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.rss = append(r.rss, addrs...)

	return nil
}

// Next 获取下一个服务器地址（轮询逻辑）
// 如果没有服务器地址，返回空字符串
func (r *RoundRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	lens := len(r.rss)
	if r.curIndex >= lens {
		r.curIndex = 0
	}

	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}

// Get 获取服务器地址（包装了Next方法）
// 参数：未使用（保留接口一致性）
// 返回值：服务器地址和错误（这里总是返回nil错误）
func (r *RoundRobinBalance) Get() (string, error) {
	return r.Next(), nil
}
