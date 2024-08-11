package channel

import (
	"sync"
)

// MutexWrap 锁代理，支持禁用
type MutexWrap struct {
	lock     sync.Mutex
	disabled bool // 是否禁用，禁用后锁失效
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}
