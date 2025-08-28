package middleware

import (
	"log"
	"net/http"
)

// ErrorHandlingMiddleware 优雅地恢复并向客户端发送一条清晰的错误消息
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 捕获panic，记录错误，向客户端发送一个 500 错误。
				log.Printf("Error  occurred:  %v", err)
				http.Error(w, "Internal  Server  Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
