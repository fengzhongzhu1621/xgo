package backoff

import (
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

// DoTempDelay 实现了指数退避策略，用于处理临时性错误的重试延迟
// 参数 tempDelay: 当前的延迟时间，如果是第一次重试则为0
// 返回值: 计算出的新延迟时间，用于下一次重试
func DoTempDelay(tempDelay time.Duration) time.Duration {
	// 判断是否为第一次重试
	if tempDelay == 0 {
		// 第一次重试，设置初始延迟为5毫秒
		tempDelay = 5 * time.Millisecond
	} else {
		// 非第一次重试，将延迟时间翻倍（指数增长核心）
		tempDelay *= 2
	}

	// 设置最大延迟上限为1秒，防止延迟时间无限增长
	if max := 1 * time.Second; tempDelay > max {
		tempDelay = max
	}

	// 让当前goroutine休眠指定的延迟时间
	time.Sleep(tempDelay)

	// 返回新的延迟时间，供下一次重试使用
	return tempDelay
}
