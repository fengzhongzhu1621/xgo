// Package throttle implement the throttle for retry/hedging request.
package throttle

import (
	"fmt"
	"math"
	"sync/atomic"
)

const (
	// MaximumTokens defines the up limit of tokens for TokenBucket.
	MaximumTokens = 1000 // 令牌桶容量上限
)

// TokenBucket defines a throttle based on token bucket.
// TokenBucket's methods may be called by multiple goroutines simultaneously.
type TokenBucket struct {
	tokens     uint64  // 当前令牌数（原子操作，存储为float64的二进制形式）
	maxTokens  float64 // 令牌桶最大容量
	threshold  float64 // 触发限流的阈值（默认maxTokens/2）
	tokenRatio float64 // 成功请求时的令牌奖励比例
}

// NewTokenBucket create a new TokenBucket.
// 令牌桶以固定速率生成令牌，请求消耗令牌。当令牌不足时，请求被限制。
// 通过OnSuccess和OnFailure方法根据请求结果增减令牌，实现自适应限流。
// 自适应限流：根据请求结果动态调整令牌数（成功奖励、失败惩罚），平衡系统负载
// //go:nosplit：禁止协程调度抢占，确保高性能（适用于高频调用的短函数）
func NewTokenBucket(maxTokens, tokenRatio float64) (*TokenBucket, error) {
	if maxTokens > MaximumTokens {
		return nil, fmt.Errorf("expect tokens less or equal to %d, got %f", MaximumTokens, maxTokens)
	}

	if maxTokens <= 0 {
		return nil, fmt.Errorf("expect positive tokens, got %f", maxTokens)
	}

	if tokenRatio <= 0 {
		return nil, fmt.Errorf("expect positive taken ratio, got %f", tokenRatio)
	}

	return &TokenBucket{
		tokens:     math.Float64bits(maxTokens), // // 初始化为满桶
		maxTokens:  maxTokens,
		threshold:  maxTokens / 2, // 动态计算的阈值（默认半满），用于快速判断是否允许请求，避免频繁计算
		tokenRatio: tokenRatio,    // 成功请求时的令牌奖励比例
	}, nil
}

// Allow whether a new request could be issued.
// 检查当前令牌是否足够发起新请求。当前令牌数超过阈值时允许请求，否则触发限流。
// 动态阈值 threshold 分离限流判断与令牌计算，降低 Allow() 的调用开销
// 失败水位下降
// 成功水位上升
//
//go:nosplit
func (tb *TokenBucket) Allow() bool {
	return math.Float64frombits(atomic.LoadUint64(&tb.tokens)) > tb.threshold
}

// OnSuccess increase tokens in bucket by token ratio, but not greater than maxTokens.
//
//go:nosplit
func (tb *TokenBucket) OnSuccess() {
	for {
		tokens := math.Float64frombits(atomic.LoadUint64(&tb.tokens))
		if tokens == tb.maxTokens {
			// 桶已满，无需增加
			return
		}

		// 增加令牌（按 tokenRatio 比例增加令牌，但不超过 maxTokens）
		newTokens := tokens + tb.tokenRatio
		if newTokens > tb.maxTokens {
			newTokens = tb.maxTokens
		}

		// 比较并交换，解决并发冲突
		if atomic.CompareAndSwapUint64(
			&tb.tokens,
			math.Float64bits(tokens),
			math.Float64bits(newTokens),
		) {
			break
		}
	}
}

// OnFailure decrease tokens in bucket by 1, but not less than 0.
//
//go:nosplit
func (tb *TokenBucket) OnFailure() {
	for {
		tokens := math.Float64frombits(atomic.LoadUint64(&tb.tokens))
		if tokens == 0 {
			return
		}

		// 失败时固定减少1个令牌（下限为0），快速响应服务降级
		// 与熔断协同：类似SRE过载保护，失败率高时自动限流
		newTokens := tokens - 1
		if newTokens < 0 {
			newTokens = 0
		}

		if atomic.CompareAndSwapUint64(
			&tb.tokens,
			math.Float64bits(tokens),
			math.Float64bits(newTokens),
		) {
			break
		}
	}
}
