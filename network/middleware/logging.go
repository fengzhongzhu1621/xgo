package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware 跟踪错误和监控 API 的行为，用于记录传入的 HTTP 请求
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 记录日志
		log.Printf("Http request Started  %s  %s", r.Method, r.URL.Path)

		// 调用链中的下一个处理程序
		next.ServeHTTP(w, r)

		// 记录所用时间
		log.Printf("Http request completed in  %v", time.Since(start))
	})
}
