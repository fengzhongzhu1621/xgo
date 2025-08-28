package middleware

import (
	"net/http"

	"github.com/fengzhongzhu1621/xgo/collections/flowctrl/ratelimiter"
)

func IPRateLimitMiddleware(limiter *ratelimiter.IPRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := limiter.GetLimiter(r.RemoteAddr)
			if !limiter.Allow() {
				// 超过限制
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
