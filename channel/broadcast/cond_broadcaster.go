package broadcast

import (
	"sync"
)

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
