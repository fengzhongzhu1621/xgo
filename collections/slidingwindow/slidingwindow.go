package slidingwindow

import (
	"sync"
	"time"
)

// SlidingWindow sliding window consists two windows `curr` and `prev`,
// the window is advanced when recording events.
// 维护两个时间窗口（curr和prev），通过滑动机制实现时间区间内的精确计数
type SlidingWindow struct {
	size time.Duration // 单个窗口的时间长度（如1秒）
	mu   sync.Mutex    // 互斥锁保证并发安全

	curr *window // 当前活跃窗口
	prev *window // 上一个窗口（用于平滑计数）
}

// NewSlidingWindow creates a new slidingwindow
func NewSlidingWindow(size time.Duration) *SlidingWindow {
	currWin := NewLocalWindow() // 	初始化当前窗口

	// The previous window is static (i.e. no add changes will happen within it),
	// so we always create it as an instance of window.
	//
	// In this way, the whole limiter, despite containing two windows, now only
	// consumes at most one goroutine for the possible sync behaviour within
	// the current window.
	prevWin := NewLocalWindow() // 初始化上一个窗口

	return &SlidingWindow{
		size: size,
		curr: currWin,
		prev: prevWin,
	}
}

// Size returns the time duration of one window size. Note that the size
// is defined to be read-only, if you need to change the size,
// create a new limiter with a new size instead.
// 获取窗口大小
func (sw *SlidingWindow) Size() time.Duration {
	return sw.size
}

// Record report whether a event may happen at time now.
// 默认记录1个事件
func (sw *SlidingWindow) Record() {
	sw.RecordN(time.Now(), 1)
}

// RecordN reports whether n events may happen at time now.
// 记录事件
func (sw *SlidingWindow) RecordN(now time.Time, n int64) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.advance(now)     // 推进窗口状态
	sw.curr.AddCount(n) // 增加当前窗口计数
}

// Count counts new window size. 计算当前窗口的加权计数
// 通过加权计算，避免了固定窗口算法在边界处的计数突变
func (sw *SlidingWindow) Count() int64 {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()

	// 推进窗口状态
	sw.advance(now)

	// 计算当前窗口已过去的时间(elapsed)
	elapsed := now.Sub(sw.curr.Start())

	// 计算权重(weight)，表示上一个窗口对当前计数的贡献比例
	// 即根据上一个窗口事件的总数，预估下一个窗口在未来的剩余时间内的事件数量
	// weight = (窗口总大小 - 已过去时间) / 窗口总大小
	weight := float64(sw.size-elapsed) / float64(sw.size)

	// 返回加权后的总计数
	//
	// weight：表示前一个窗口对当前计数的贡献比例
	// sw.prev.Count()：前一个完整窗口的计数
	// sw.curr.Count()：当前窗口的计数
	// count = 当前窗口已过去的时间内的事件总数 + 窗口剩余时间的事件总数估计（根据上一个窗口的事件总数进行预估）
	count := int64(weight*float64(sw.prev.Count())) + sw.curr.Count()

	return count
}

// advance updates the current/previous windows resulting from the passage of time.
func (sw *SlidingWindow) advance(now time.Time) {
	// Calculate the start boundary of the expected current-window.
	// 计算新窗口起始时间
	// 将当前时间 now 对齐到最近的、小于或等于 now 的 sw.size 时间单位的整倍数时刻
	newCurrStart := now.Truncate(sw.size)
	// 当newCurrStart与当前窗口起始时间的差值超过size时触发滑动
	diffSize := newCurrStart.Sub(sw.curr.Start()) / sw.size

	// 如果diffSize == 0，仍在当前窗口，不做任何操作
	// 如果diffSize == 1，滑动一个窗口，将当前窗口计数转移到前一个窗口
	// 如果diffSize > 1，说明跳过了多个窗口，重置两个窗口计数为0

	// Fast path, the same window
	if diffSize == 0 {
		// 仍在当前窗口
		return
	}

	// 触发滑动
	// Slow path, the current-window is at least one-window-size behind the expected one.

	// The new current-window always has zero count.
	// 重置当前窗口，计数归零
	sw.curr.Reset(newCurrStart, 0)

	// 如果只滑动了一个窗口大小(diffSize == 1)，则将当前窗口计数继承给上一个窗口
	// 否则将上一个窗口计数归零
	// reset previous window
	newPrevCount := int64(0)
	if diffSize == 1 {
		// The new previous-window will overlap with the old current-window,
		// so it inherits the count.
		//
		// Note that the count here may be not accurate, since it is only a
		// SNAPSHOT of the current-window's count, which in itself tends to
		// be inaccurate due to the asynchronous nature of the sync behaviour.
		// 时间轴: |----|----|----|----|----| (每个|代表一个窗口边界)
		//        t0   t1   t2   t3   t4   t5
		//
		// 初始状态:
		// curr: [t0, t1)
		// prev: [t-1, t0)
		//
		// 时间推进到t1.5:
		// newCurrStart = t1
		// diffSize = (t1 - t0)/size = 1
		//
		// 滑动后状态:
		// curr: [t1, t2), count=0
		// prev: [t0, t1), count=原curr.count
		newPrevCount = sw.curr.Count()
	} else {
		// 当 diffSize > 1 时，说明当前时间已经跨越了多个完整的窗口周期。这种情况通常发生在：
		//
		// 系统负载突然降低，导致长时间没有请求到达
		// 系统时钟被手动调整
		// 程序暂停执行后恢复（如GC暂停）
		//
		// 时间轴: |----|----|----|----|----| (每个|代表一个窗口边界，窗口大小为1单位时间)
		//        t0   t1   t2   t3   t4   t5
		//
		// 初始状态:
		// curr: [t0, t1)
		// prev: [t-1, t0)
		//
		// 时间推进到t3.5:
		// newCurrStart = t3
		// diffSize = (t3 - t0)/size = 3
		//
		// 滑动后状态:
		// curr: [t3, t4), count=0
		// prev: [t2, t3), count=0 (因为diffSize > 1，不继承任何计数)
	}

	sw.prev.Reset(newCurrStart.Add(-sw.size), newPrevCount)
}
