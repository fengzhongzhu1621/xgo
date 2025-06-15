package lock

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
)

// 可重入锁又称为递归锁，是指在同一个线程在外层方法获取锁的时候，在进入该线程的内层方法时会自动获取锁，不会因为之前已经获取过还没释放再次加锁导致死锁
// Go 里面的 Mutex 不是可重入的锁。Mutex 的实现中没有记录哪个 goroutine 拥有这把锁。
// 理论上，任何 goroutine 都可以随意地 Unlock 这把锁，所以没办法计算重入条件，并且Mutex 重复Lock会导致死锁。

// ReentrantMutex 是一个可重入的互斥锁
type ReentrantMutex struct {
	mu        sync.Mutex
	owner     int64
	recursion int32
}

// Lock 尝试获取锁
func (m *ReentrantMutex) Lock() {
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.mu.Lock()
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

// Unlock 释放锁
func (m *ReentrantMutex) Unlock() {
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	m.recursion--
	if m.recursion > 0 {
		return
	}
	atomic.StoreInt64(&m.owner, 0)
	m.mu.Unlock()
}
