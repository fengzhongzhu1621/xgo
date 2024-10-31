package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stretchr/testify/assert"
)

func TestWaitGroupTimeout_no_timeout(t *testing.T) {
	// sync.WaitGroup 是 Go 语言中的一个同步原语，用于等待一组协程（goroutine）完成执行。
	// 它提供了一种简单的方式来同步多个并发执行的协程，确保在所有协程完成之前，主协程不会退出。
	wg := &sync.WaitGroup{}

	timeouted := WaitGroupTimeout(wg, time.Millisecond*100)
	assert.False(t, timeouted)
}

func TestWaitGroupTimeout_timeout(t *testing.T) {
	wg := &sync.WaitGroup{}
	// 增加等待组的计数器。delta 参数表示要增加的值。通常在启动一个新的协程之前调用此方法。
	wg.Add(1)

	// 因为没有执行 wg.Done() // 减少等待组计数器，wg 会超时退出
	timeouted := WaitGroupTimeout(wg, time.Millisecond*100)
	assert.True(t, timeouted)
}

func TestSizedwaitgroup(t *testing.T) {
	var wg sync.WaitGroup
	// 创建一个容量为 3 的 SizedWaitGroup, 最多只能有 3 个协程同时运行
	swg := sizedwaitgroup.New(3)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			swg.Add()        // 增加 SizedWaitGroup 的计数器
			defer swg.Done() // 减少 SizedWaitGroup 的计数器

			// 当一个协程完成任务并调用 swg.Done() 时，SizedWaitGroup 的计数器会减少，从而允许其他等待中的协程开始执行。
			// 模拟一个耗时任务
			time.Sleep(time.Second)
			fmt.Printf("Task %d completed\n", i)
		}(i)
	}

	// 还需要使用标准库中的 sync.WaitGroup 来确保主协程等待所有任务完成。
	// 这是因为 sizedwaitgroup 只负责限制并发数量，而不负责等待所有任务完成。
	wg.Wait() // 等待所有协程完成
}
