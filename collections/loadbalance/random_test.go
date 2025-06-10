package loadbalance

import (
	"fmt"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRandomBalance(t *testing.T) {
	lb := &RandomBalance{}

	// 添加服务器地址
	if err := lb.Add("192.168.1.1:8080", "192.168.1.2:8080", "192.168.1.3:8080"); err != nil {
		fmt.Println("Error adding servers:", err)
		return
	}

	var g errgroup.Group

	// 模拟并发获取服务器
	for i := 0; i < 5; i++ {
		g.Go(func() error {
			server, err := lb.Get()
			if err != nil {
				fmt.Println("Error getting server:", err)
			} else {
				fmt.Println("Selected server:", server)
			}

			return nil
		})
	}

	// 等待一段时间以观察输出
	g.Wait()
}
