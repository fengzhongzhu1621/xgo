package atomicutils

import (
	"sync/atomic"
)

// Counter 自增计数器.
type Counter struct {
	val int32
}

func (c *Counter) Increment() {
	// 将 val 增加 1。
	atomic.AddInt32(&c.val, 1)
}

func (c *Counter) Value() int32 {
	// 原子地读取 val 的当前值。
	return atomic.LoadInt32(&c.val)
}

func (c *Counter) Reset() {
	// 原子地将 val 设置为 0。
	atomic.StoreInt32(&c.val, 0)
}
