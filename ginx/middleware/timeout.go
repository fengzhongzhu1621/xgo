package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout" // 提供超时中间件
	"github.com/gin-gonic/gin"       // Gin Web 框架
)

func timeoutResponse(c *gin.Context) {
	c.String(http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout))
}

func TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(500*time.Millisecond), // 设置超时时间为 500ms
		timeout.WithHandler(func(c *gin.Context) { // 正常的请求处理逻辑（这里只是调用 c.Next()）
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse), // 超时时调用的响应函数
	)
}
