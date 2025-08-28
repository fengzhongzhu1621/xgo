package middleware

import "net/http"

// ValidationMiddleware 一个验证中间件示例
func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") == "" {
			http.Error(w, "Missing API Key", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
