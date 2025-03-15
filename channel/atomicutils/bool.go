package atomicutils

import (
	"sync/atomic"
)

// AtomicBool 使用标准库的 atomic.Bool
type AtomicBool struct {
	b atomic.Bool
}

// NewAtomicBool 创建一个新的 AtomicBool 实例，并根据 yes 初始化
func NewAtomicBool(yes bool) *AtomicBool {
	ab := &AtomicBool{}
	ab.b.Store(yes)
	return ab
}

// SetIfNotSet 如果当前值为 false，则设置为 true，并返回 true；否则返回 false
func (ab *AtomicBool) SetIfNotSet() bool {
	for {
		old := ab.b.Load()
		if old {
			return false
		}
		// ​比较：检查内存位置的当前值是否等于预期值。
		// ​交换：如果相等，则将该位置的值更新为新值；否则，不进行任何操作。
		if ab.b.CompareAndSwap(false, true) {
			return true
		}
		// CAS 失败，重试
	}
}

// Set 将值设置为 true
func (ab *AtomicBool) Set() {
	ab.b.Store(true)
}

// Unset 将值设置为 false
func (ab *AtomicBool) Unset() {
	ab.b.Store(false)
}

// IsSet 检查当前值是否为 true
func (ab *AtomicBool) IsSet() bool {
	return ab.b.Load()
}

// SetTo 根据 yes 设置值为 true 或 false
func (ab *AtomicBool) SetTo(yes bool) {
	ab.b.Store(yes)
}
