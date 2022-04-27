package atomicutils

import (
	"sync/atomic"
)

// Counter 自增计数器
type Counter struct {
	val int32
}

func (c *Counter) Increment() {
	atomic.AddInt32(&c.val, 1)
}

func (c *Counter) Value() int32 {
	return atomic.LoadInt32(&c.val)
}

func (c *Counter) Reset() {
	atomic.StoreInt32(&c.val, 0)
}
