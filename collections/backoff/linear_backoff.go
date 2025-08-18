package backoff

import (
	"errors"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// ////////////////////////////////////////////////////////////////////////////////
type linearBackoff time.Duration

// NextBackoff 下一次重试前的时间间隔，实现了一个恒定的重试间隔。
func (r linearBackoff) NextBackoff() time.Duration {
	return time.Duration(r)
}

func NewLinearBackoff(backoff time.Duration) IRetryStrategy {
	return linearBackoff(backoff)
}

// ////////////////////////////////////////////////////////////////////////////////
// linearBackoff has lowest priority in all kinds of backoff.
// bfs：一个 time.Duration 类型的切片，用于存储预设的线性退避时间间隔序列。例如，[1s, 2s, 3s] 表示第一次重试等待1秒，第二次2秒，第三次3秒
type linearBackoffs struct {
	bfs []time.Duration // 存储线性递增的退避时间序列
}

// NewLinearBackoffs create a new linear backoff. Empty bfs will cause an error.
func NewLinearBackoffs(bfs ...time.Duration) (*linearBackoffs, error) {
	if len(bfs) == 0 {
		return nil, errors.New("linear backoff list must not be empty")
	}
	return &linearBackoffs{bfs: bfs}, nil
}

// backoff is randomly distributed in [0, bfs[min(len(bfs)-1, attempt)]].
func (bf *linearBackoffs) Backoff(attempt int) (delay time.Duration) {
	defer func() {
		// 对最终计算的 delay 乘以一个 [0, 1) 的随机数，使实际延迟时间在 [0, delay) 之间均匀分布。这避免了多个客户端同时重试导致的“惊群效应”
		delay = time.Duration(rand.Float64() * float64(delay))
	}()

	if attempt <= 0 {
		return 0
	}

	// 返回 bfs 中第 attempt-1 个元素（切片索引从0开始）。例如，第一次重试（attempt=1）返回 bfs[0]
	l := len(bf.bfs)
	if attempt <= l {
		return bf.bfs[attempt-1]
	}

	// 返回 bfs 的最后一个元素，表示达到最大退避时间后不再递增
	return bf.bfs[l-1]
}
