package ginx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type (
	ResponseAssertFunc func(Response) error
	JSONAssertFunc     func(map[string]interface{}) error
)

// SetupRouter 创建了一个新的Gin路由引擎实例
func SetupRouter() *gin.Engine {
	// 创建了一个新的Gin路由引擎实例
	r := gin.New()
	// 设置了Gin的运行模式为发布模式。在发布模式下，Gin会关闭调试信息，以提高性能。
	// 如果你想在开发过程中启用调试信息，可以使用gin.DebugMode
	gin.SetMode(gin.ReleaseMode)
	// 添加了一个中间件，该中间件会在每个请求处理过程中捕获任何可能发生的panic，并恢复程序的正常执行。
	// 这对于防止因为未处理的异常导致整个服务崩溃非常有用。
	r.Use(gin.Recovery())

	return r
}

func SetupDebugModeRouter() *gin.Engine {
	r := gin.New()
	gin.SetMode(gin.DebugMode)
	r.Use(gin.Recovery())

	return r
}

func NewJSONAssertFunc(
	t assert.TestingT,
	assertFunc JSONAssertFunc,
) func(res *http.Response, req *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		// 读取响应结果
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err, "read body from response fail")

		defer res.Body.Close()

		// 将 json 字符串转换为结构体
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		assert.NoError(t, err, "unmarshal string to json fail")

		return assertFunc(data)
	}
}

func NewResponseAssertFunc(
	t *testing.T,
	responseFunc ResponseAssertFunc,
) func(res *http.Response, req *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err, "read body from response fail")

		defer res.Body.Close()

		var data Response

		err = json.Unmarshal(body, &data)
		assert.NoError(t, err, "unmarshal string to response fail")

		return responseFunc(data)
	}
}

// CreateTestContextWithDefaultRequest ...
func CreateTestContextWithDefaultRequest(w *httptest.ResponseRecorder) *gin.Context {
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest("POST", "/", new(bytes.Buffer))
	return ctx
}
