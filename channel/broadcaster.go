package channel

import (
	"context"
	"sync"
	"sync/atomic"
)

type ChannelBroadcaster struct {
	signal chan struct{} // 定义广播信号
}

func NewChannelBroadcaster() *ChannelBroadcaster {
	return &ChannelBroadcaster{
		signal: make(chan struct{}),
	}
}

// Go 等待接受广播
func (b *ChannelBroadcaster) Go(fn func()) {
	go func() {
		<-b.signal
		fn()
	}()
}

// Broadcast 发送广播
func (b *ChannelBroadcaster) Broadcast() {
	close(b.signal)
}

type CondBroadcaster struct {
	mu       *sync.Mutex
	cond     *sync.Cond
	signaled bool
}

func NewCondBroadcaster() *CondBroadcaster {
	var mu sync.Mutex
	return &CondBroadcaster{
		mu:       &mu,
		cond:     sync.NewCond(&mu),
		signaled: false,
	}
}

func (b *CondBroadcaster) Go(fn func()) {
	go func() {
		b.cond.L.Lock()
		defer b.cond.L.Unlock()

		for !b.signaled {
			// 协程启动后在这里等待广播信号
			b.cond.Wait()
		}
		fn()
	}()
}

func (b *CondBroadcaster) Broadcast() {
	b.cond.L.Lock()
	b.signaled = true
	b.cond.L.Unlock()

	b.cond.Broadcast()
}

type ContextBroadcaster struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewContextBroadcaster() *ContextBroadcaster {
	ctx, cancel := context.WithCancel(context.Background())
	return &ContextBroadcaster{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (b *ContextBroadcaster) Go(fn func()) {
	go func() {
		<-b.ctx.Done()
		fn()
	}()
}

func (b *ContextBroadcaster) Broadcast() {
	b.cancel()
}

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

type RWMutexBroadcaster struct {
	mu *sync.RWMutex
}

func NewRWMutexBroadcaster() *RWMutexBroadcaster {
	var mu sync.RWMutex
	// 添加读写锁
	mu.Lock()

	return &RWMutexBroadcaster{mu: &mu}
}

func (b *RWMutexBroadcaster) Go(fn func()) {
	go func() {
		// 读锁定，等待解锁
		b.mu.RLock()
		defer b.mu.RUnlock()

		fn()
	}()
}

func (b *RWMutexBroadcaster) Broadcast() {
	b.mu.Unlock()
}
