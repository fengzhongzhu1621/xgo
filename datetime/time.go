package datetime

import (
	"context"
	"time"
)

// Sleep 休眠指定时间.
func Sleep(ctx context.Context, dur time.Duration) error {
	// 创建一个定时器
	t := time.NewTimer(dur)
	defer t.Stop()

	select {
	case <-t.C: // 在指定间隔后从定时器获取下一个调度时间，如果还未到时间则阻塞
		return nil
	case <-ctx.Done(): // 返回一个 Channel，Channel会在当前工作完成或者上下文被取消后关闭，多次调用Done方法会返回同一个Channel
		return ctx.Err()
	}
}

// SecondToTime 把时间戳转换为time.Time格式.
func SecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

// TodayStartTimestamp 返回今天开始时间的 Unix 时间戳（自 1970 年 1月1日 00:00:00 UTC 起经过的秒数）。
func TodayStartTimestamp() int64 {
	now := time.Now()
	year, month, day := now.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	return start.Unix()
}
