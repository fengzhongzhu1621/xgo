package xerror

import (
	"sync"
	"time"
)

// ErrFrequencyInterface 跟踪错误的发生频率，并在一定时间窗口内判断某个错误是否总是出现。
type ErrFrequencyInterface interface {
	// IsErrAlwaysAppear 判断给定的错误是否总是出现。
	IsErrAlwaysAppear(err error) bool

	// Release 释放错误，重置状态。
	Release()
}

type errFrequency struct {
	mu      sync.Mutex
	err     error         // 当前跟踪的错误
	window  time.Duration // 错误检测窗口
	endTime time.Time     // 错误跟踪的结束时间戳
}

// NewErrFrequency 创建一个新的 errFrequency 实例，并设置初始错误和结束时间。
func NewErrFrequency(err error, duration time.Duration) ErrFrequencyInterface {
	return &errFrequency{
		err:     err,
		window:  duration,
		endTime: time.Now().Add(duration),
	}
}

// IsErrAlwaysAppear 判断给定的错误是否总是出现。
func (e *errFrequency) IsErrAlwaysAppear(err error) bool {
	if err == nil {
		return false
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// 如果当前跟踪的错误与传入的错误相同，并且当前时间已经超过结束时间，返回 true
	if e.err != nil && e.err.Error() == err.Error() {
		return time.Now().After(e.endTime)
	}

	e.err = err
	e.endTime = time.Now().Add(e.window)
	return false
}

// Release 重置当前跟踪的错误
func (e *errFrequency) Release() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.err = nil
	e.endTime = time.Time{}
}
