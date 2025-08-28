package counter

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// 具体的实现参考： /channel/atomicutils/*
func TestAtomicCounter(t *testing.T) {
	var counter int64 // 声明一个int64类型的变量counter，用于计数

	var wg sync.WaitGroup // 声明一个WaitGroup，用于等待所有goroutine完成
	numWorkers := 10      // 定义工作goroutine的数量为10

	wg.Add(numWorkers) // 将WaitGroup的计数器设置为numWorkers，表示有10个goroutine需要等待
	for i := 0; i < numWorkers; i++ {
		go func() {
			atomic.AddInt64(&counter, 1) // 原子性地对counter加1，保证并发安全
			wg.Done()                    // 表示一个goroutine已完成，WaitGroup计数器减1
		}()
	}

	wg.Wait()                              // 主线程阻塞，直到WaitGroup的计数器归零，即所有goroutine都已完成
	fmt.Println("Counter value:", counter) // 打印最终的计数器值
}
