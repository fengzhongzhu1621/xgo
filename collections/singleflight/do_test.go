package singleflight

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"
)

var g singleflight.Group

// 模拟一个耗时的计算
func expensiveComputation(index int) (interface{}, error) {
	fmt.Println("Performing expensive computation...")
	time.Sleep(2 * time.Second)
	return index, nil
}

func TestDo(t *testing.T) {
	var wg sync.WaitGroup

	// 定义一个 key
	gKey := "key"

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 每个 goroutine 都发起耗时计算请求
			result, err, _ := g.Do(gKey, func() (interface{}, error) {
				return expensiveComputation(i)
			})

			if err != nil {
				fmt.Printf("Goroutine %d: Error: %v\n", id, err)
				return
			}

			fmt.Printf("Goroutine %d: Got result: %v\n", id, result)
		}(i)
	}

	wg.Wait()

	// 只有一个请求会真正执行计算，其他请求会共享结果。
	// Performing expensive computation...
	// Goroutine 1: Got result: 2
	// Goroutine 2: Got result: 2
	// Goroutine 0: Got result: 2
}
