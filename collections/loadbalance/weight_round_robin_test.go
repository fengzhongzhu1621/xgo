package loadbalance

import (
	"fmt"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestWeightRoundRobinBalance(t *testing.T) {
	lb := &WeightRoundRobinBalance{}

	// 添加服务器地址
	lb.Add("192.168.1.1:8080", 100)
	lb.Add("192.168.1.2:8080", 100)
	lb.Add("192.168.1.3:8080", 100)

	// 模拟并发获取服务器
	var g errgroup.Group
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

	g.Wait()
}
