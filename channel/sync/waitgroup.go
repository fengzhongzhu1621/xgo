package sync

import (
	"sync"
	"time"
)

// WaitGroupTimeout adds timeout feature for sync.WaitGroup.Wait().
// It returns true, when timeouted.
func WaitGroupTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	wgClosed := make(chan struct{}, 1)
	go func() {
		// 等待所有协程完成
		wg.Wait()
		wgClosed <- struct{}{}
	}()

	select {
	case <-wgClosed:
		// 没有超时，返回 false
		return false
	case <-time.After(timeout):
		// 超时返回 true
		return true
	}
}
