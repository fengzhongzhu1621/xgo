package loadbalance

import (
	"fmt"
	"testing"
	"time"
)

func TestRandomBalance(t *testing.T) {
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
