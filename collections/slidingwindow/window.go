package slidingwindow

import (
	"time"
)

// window represents a window that ignores sync behavior entirely
// and only stores counters in memory.
// 实现了一个简单的本地滑动窗口计数器，用于统计时间窗口内的事件数量（如请求次数）
type window struct {
	// The start boundary (timestamp in nanoseconds) of the window.
	// [start, start + size)
	start int64 // 窗口的起始时间（纳秒精度），标记窗口的时间范围 [start, start + size)。

	// The total count of events happened in the window.
	count int64 // 当前窗口内累计的事件数量（如请求次数）
}

// NewLocalWindow 返回一个空窗口（start和count初始化为零值）
func NewLocalWindow() *window {
	return &window{}
}

// Start 返回窗口的起始时间（time.Time 类型），用于外部查询窗口的时间范围
func (w *window) Start() time.Time {
	return time.Unix(0, w.start)
}

// Count 获取事件计数
func (w *window) Count() int64 {
	return w.count
}

// AddCount 增加事件计数
func (w *window) AddCount(n int64) {
	w.count += n // 原子操作增加计数（实际需外部加锁）
}

// Reset 重置窗口状态
func (w *window) Reset(s time.Time, c int64) {
	w.start = s.UnixNano() // 更新起始时间（转为纳秒）
	w.count = c            // 重置计数
}
