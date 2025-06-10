package loadbalance

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
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
	if len(r.rss) == 0 {
		return "", errors.New("no servers available")
	}
	index := rand.Intn(len(r.rss))
	return r.rss[index], nil
}

// Get 根据调用获取一个服务器地址
func (r *RandomBalance) Get() (string, error) {
	return r.Next()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	lb := &RandomBalance{}

	// 添加服务器地址
	if err := lb.Add("192.168.1.1:8080", "192.168.1.2:8080", "192.168.1.3:8080"); err != nil {
		fmt.Println("Error adding servers:", err)
		return
	}

	// 模拟并发获取服务器
	for i := 0; i < 5; i++ {
		go func() {
			server, err := lb.Get()
			if err != nil {
				fmt.Println("Error getting server:", err)
			} else {
				fmt.Println("Selected server:", server)
			}
		}()
	}

	// 等待一段时间以观察输出
	time.Sleep(2 * time.Second)
}
