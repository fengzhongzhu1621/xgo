package channel

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func cal() {
	for i := 0; i < 1000000; i++ {
		runtime.Gosched()
	}
}

// 计算Goroutine 切换的开销
func TestGosched(t *testing.T) {
	runtime.GOMAXPROCS(1) // 设置使用1个CPU核心

	currentTime := time.Now()
	fmt.Println(currentTime)

	go cal() // 启动一个goroutine执行cal函数

	for i := 0; i < 1000000; i++ {
		runtime.Gosched()
	}

	// 87ns
	fmt.Println(time.Since(currentTime) / 2000000)
}

// 计算协程使用的内存开销
func TestGetGoroutineMemConsume(t *testing.T) {
	var c chan int
	var wg sync.WaitGroup
	const goroutineNum = 1000000

	memConsumed := func() uint64 {
		runtime.GC() // 触发垃圾回收，排除对象影响
		var memStat runtime.MemStats
		runtime.ReadMemStats(&memStat)
		return memStat.Sys
	}

	noop := func() {
		wg.Done()
		<-c // 防止 goroutine 退出，内存被释放
	}

	wg.Add(goroutineNum)
	before := memConsumed() // 获取创建 goroutine 前的内存
	for i := 0; i < goroutineNum; i++ {
		go noop()
	}
	wg.Wait()
	after := memConsumed() // 获取创建 goroutine 后的内存

	fmt.Println(runtime.NumGoroutine())
	fmt.Printf("%.3f KB bytes\n", float64(after-before)/goroutineNum/1024)
}
