package channel

import (
	"sync"
	"time"
)

// WaitGroupTimeout adds timeout feature for sync.WaitGroup.Wait().
// It returns true, when timeouted.
func WaitGroupTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	// 异步等待所有协程完成
	wgClosed := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		wgClosed <- struct{}{}
	}()

	// 等待任务执行完毕通知
	select {
	case <-wgClosed:
		// 所有任务执行完成未超时，返回 false
		return false
	case <-time.After(timeout):
		// 超时返回 true
		return true
	}
}
