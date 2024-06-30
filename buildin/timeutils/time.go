package timeutils

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// RetryBackoff 计算重试的间隔
// retry: 重试次数
// minBackoff: 最小时间间隔
// maxBackoff: 最大时间间隔.
func RetryBackoff(retry int, minBackoff, maxBackoff time.Duration) (time.Duration, error) {
	if retry < 0 {
		return 0, fmt.Errorf("not reached")
	}
	if minBackoff == 0 {
		return 0, nil
	}
	// 根据重试次数计算随机范围
	d := minBackoff << uint(retry)
	if d < minBackoff {
		return maxBackoff, nil
	}
	// 随机范围在 [minBackoff, minBackoff << uint(retry)] 之间
	d = minBackoff + time.Duration(rand.Int63n(int64(d)))

	if d > maxBackoff || d < minBackoff {
		d = maxBackoff
	}

	return d, nil
}

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

// SecondToTime 把事件戳转换为time.Time格式.
func SecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}
