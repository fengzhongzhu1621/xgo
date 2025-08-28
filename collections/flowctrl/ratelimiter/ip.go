package ratelimiter

import (
	"sync"

	"golang.org/x/time/rate"
)

// IPRateLimiter 基于IP的速率限制器（令牌桶算法）
type IPRateLimiter struct {
	ips map[string]*rate.Limiter // 基于IP的限制，存储每个IP对应的限流器，考虑使用sync.Map代替map+RWMutex，特别是当IP数量很大时
	mu  *sync.RWMutex
	// 表示事件的速率限制。Limit 可能是一个自定义类型，通常表示每秒允许的事件数。
	// 每秒生成的令牌速率（即速率限制）令牌/秒（tokens per second, TPS）
	// rate.Limit(1) → 1 TPS（每秒 1 个请求）
	// rate.Limit(5) → 5 TPS（每秒 5 个请求）
	// rate.Inf → 无限速率（不限制）
	r rate.Limit
	// 表示允许的最大突发量，即短时间内允许的最大事件数。/ 令牌桶的容量（即突发量）
	// 令牌桶初始时是满的（即一开始有 b 个令牌）。
	// 如果请求速率低于 r，多余的令牌会累积在桶中，最多不超过 b 个。
	// 如果桶中没有令牌，新的请求会被阻塞（或返回 false，取决于调用方式）。
	// b = 3 → 允许最多 3 个突发请求（即使当前没有新令牌生成）。
	// b = 0 → 不允许任何突发（必须严格按照 r 的速率请求）。
	// b = math.MaxInt32 → 理论上允许无限突发（但实际受系统限制）。
	b int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b, // 表示允许的最大突发量，即短时间内允许的最大事件数。
	}
}

func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}
