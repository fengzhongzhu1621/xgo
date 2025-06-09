package channel

import (
	"fmt"
	"runtime"
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
	fmt.Println(time.Now().Sub(currentTime) / 2000000)
}
