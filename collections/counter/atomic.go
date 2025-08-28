package counter

import "sync/atomic"

// AtomicCounter 使用原子操作实现线程安全的计数器
// 使用原子操作的 AtomicCounter 在高并发场景下性能显著优于使用互斥锁的 MutexCounter
type AtomicCounter struct {
	value int64
}

// Increment 原子性地增加计数器的值
func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

// Value 原子性地返回当前计数器的值
func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}
