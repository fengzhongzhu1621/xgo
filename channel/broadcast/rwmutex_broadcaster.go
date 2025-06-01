package broadcast

import (
	"sync"
)

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
