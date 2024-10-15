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
