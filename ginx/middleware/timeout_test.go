package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTimeoutMiddleware(t *testing.T) {
	r := gin.New() // 创建一个新的 Gin 引擎（无默认中间件）

	req, _ := http.NewRequest("GET", "/slow", nil)
	w := httptest.NewRecorder()

	// 使用超时中间件
	r.Use(TimeoutMiddleware())

	// 定义一个慢速路由（模拟耗时操作）
	r.GET("/slow", func(c *gin.Context) {
		time.Sleep(800 * time.Millisecond) // 模拟耗时 800ms 的操作
		c.Status(http.StatusOK)            // 正常情况下返回 HTTP 200
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, 408, w.Code)
}

func emptySuccessResponse(c *gin.Context) {
	// 模拟业务逻辑
	time.Sleep(200 * time.Microsecond)
	c.String(http.StatusOK, "")
}

func TestTimeoutForSingleRouter(t *testing.T) {
	r := gin.New() // 创建一个新的 Gin 引擎（无默认中间件）

	req, _ := http.NewRequest("GET", "/slow", nil)
	w := httptest.NewRecorder()

	// 定义一个慢速路由（模拟耗时操作）
	r.GET("/slow", timeout.New(
		timeout.WithTimeout(100*time.Microsecond),
		timeout.WithHandler(emptySuccessResponse),
	))

	r.ServeHTTP(w, req)

	assert.Equal(t, 408, w.Code)
}
