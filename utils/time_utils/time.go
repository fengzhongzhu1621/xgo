package time_utils

import (
	"context"
	"math/rand"
	"time"
)

// 计算重试的间隔
func RetryBackoff(retry int, minBackoff, maxBackoff time.Duration) time.Duration {
	if retry < 0 {
		panic("not reached")
	}
	if minBackoff == 0 {
		return 0
	}

	d := minBackoff << uint(retry)
	if d < minBackoff {
		return maxBackoff
	}

	d = minBackoff + time.Duration(rand.Int63n(int64(d)))

	if d > maxBackoff || d < minBackoff {
		d = maxBackoff
	}

	return d
}

// 休眠指定时间
func Sleep(ctx context.Context, dur time.Duration) error {
	// 创建一个定时器
	t := time.NewTimer(dur)
	defer t.Stop()

	select {
	case <-t.C:	// 在指定间隔后从定时器获取下一个调度时间，如果还未到时间则阻塞
		return nil
	case <-ctx.Done():	// 返回一个 Channel，Channel会在当前工作完成或者上下文被取消后关闭，多次调用Done方法会返回同一个Channel
		return ctx.Err()
	}
}
