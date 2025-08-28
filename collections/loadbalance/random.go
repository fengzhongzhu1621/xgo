package loadbalance

import (
	"crypto/rand"
	"errors"
	"math/big"
	"sync"
)

// RandomBalance 随机负载均衡器
type RandomBalance struct {
	rss []string
	mu  sync.Mutex
}

// Add 向负载均衡器中添加一个或多个服务器地址
func (r *RandomBalance) Add(addrs ...string) error {
	if len(addrs) == 0 {
		return errors.New("at least one server address is required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rss = append(r.rss, addrs...)
	return nil
}

// Next 随机选择并返回一个服务器地址
func (r *RandomBalance) Next() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	total := len(r.rss)
	if total == 0 {
		return "", errors.New("no servers available")
	}
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(total)))
	return r.rss[index.Int64()], nil
}

// Get 根据调用获取一个服务器地址
func (r *RandomBalance) Get() (string, error) {
	return r.Next()
}
