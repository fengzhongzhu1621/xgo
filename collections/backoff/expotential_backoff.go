package backoff

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

// exponentialBackoff has a priority higher than linearBackoff but lower than customizedBackoff.
type ExponentialBackoff struct {
	initial    time.Duration // 初始等待时间（如100ms）
	maximum    time.Duration // 最大等待时间上限（如10秒）
	multiplier int           // 指数乘数（通常为2，表示每次等待时间翻倍）
}

// NewExponentialBackoff create a new exponentialBackoff.
func NewExponentialBackoff(initial, maximum time.Duration, multiplier int) (*ExponentialBackoff, error) {
	if initial <= 0 {
		return nil, errors.New("initial of exponential backoff must be positive")
	}

	if maximum < initial {
		return nil, errors.New("maximum of exponential backoff must be greater than initial")
	}

	if multiplier <= 0 {
		return nil, errors.New("multiplier of exponential backoff must be positive")
	}

	return &ExponentialBackoff{
		initial:    initial,
		maximum:    maximum,
		multiplier: multiplier,
	}, nil
}

// backoff returns the backoff after each attempt. The result is randomized.
func (bf *ExponentialBackoff) Backoff(attempt int) time.Duration {
	// 指数计算：math.Pow(float64(bf.multiplier), float64(attempt-1)) 计算当前尝试次数对应的指数增长倍数（如第3次尝试为multiplier^2）。
	// 时间上限限制：math.Min确保计算结果不超过maximum，避免无限等待
	ceil := math.Min(
		float64(bf.initial)*math.Pow(
			float64(bf.multiplier),
			float64(attempt-1)),
		float64(bf.maximum))
	// 随机抖动（Jitter）：rand.Float64() * ceil 在计算结果上乘以随机数（0~1之间），
	// 分散重试请求的峰值，避免“惊群效应”（多个客户端同时重试导致服务雪崩）
	return time.Duration(rand.Float64() * ceil)
}
