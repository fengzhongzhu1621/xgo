package sync

import (
	"sync"
	"testing"
	"time"

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
