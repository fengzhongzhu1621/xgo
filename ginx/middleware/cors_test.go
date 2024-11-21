package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CORS([]string{"http://1.com", "http://2.com"}))
	router.GET("/cros", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	tests := []struct {
		name               string
		origin             string
		expectedStatusCode int
	}{
		{
			"Allowed origin",
			"http://1.com",
			http.StatusOK, // 命中规则
		},
		{
			"Not allowed origin",
			"http://3.com",
			http.StatusForbidden, // 没有命中规则
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建了一个新的HTTP请求
			req, _ := http.NewRequest(http.MethodGet, "/cros", nil)
			// 设置了请求头中的Origin字段
			req.Header.Set("Origin", tt.origin)
			// 创建了一个新的httptest.ResponseRecorder实例，它实现了http.ResponseWriter接口
			// 用于记录HTTP响应，以便后续可以检查响应的状态码、头部和主体
			w := httptest.NewRecorder()
			// 将创建的HTTP请求req传递给router，并使用ResponseRecorder实例w来接收响应
			router.ServeHTTP(w, req)

			// 检查响应的状态码是否与测试用例中预期的状态码相匹配
			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
