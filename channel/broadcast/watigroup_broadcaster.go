package broadcast

import (
	"sync"
	"sync/atomic"
)

type WaitGroupBroadcaster struct {
	done    int32
	trigger sync.WaitGroup
}

func NewWaitGroupBroadcaster() *WaitGroupBroadcaster {
	b := &WaitGroupBroadcaster{}
	b.trigger.Add(1)
	return b
}

func (b *WaitGroupBroadcaster) Go(fn func()) {
	go func() {
		if atomic.LoadInt32(&b.done) == 1 {
			// 如果 done 的值为 1，表示已经广播过，那么新的 goroutine 将直接返回，不会执行 fn 函数。
			return
		}

		// 如果 done 的值为 0，表示还没有广播过，那么新的 goroutine 将调用 b.trigger.Wait() 方法等待 WaitGroup 的计数器变为 0。
		// 当计数器变为 0 时，WaitGroup 会释放等待的 goroutine，然后执行 fn 函数。
		b.trigger.Wait()
		fn()
	}()
}

func (b *WaitGroupBroadcaster) Broadcast() {
	// 尝试将 done 字段的值从 0 更改为 1
	// 如果 b.done 指向的值等于 old，则将其设置为 new，并返回 true。
	// 否则，不进行任何操作，并返回 false
	if atomic.CompareAndSwapInt32(&b.done, 0, 1) {
		// 如果成功，表示这是第一次广播，那么调用 b.trigger.Done() 方法将 WaitGroup 的计数器减 1，从而释放所有等待的 goroutine。
		b.trigger.Done()
	}
}
