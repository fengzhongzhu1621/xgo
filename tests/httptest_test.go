package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func TestHttpConn(t *testing.T) {
	// httptest.NewRequest 用于创建一个 HTTP 请求对象，通常用于测试 HTTP 服务器的处理逻辑。
	// 创建的请求对象不会发送到网络上的服务器，而是直接在内存中进行处理
	req := httptest.NewRequest("GET", "http://xxx/xx", nil)

	// 用于创建一个实现了 http.ResponseWriter 接口的对象，通常用于测试 HTTP 处理函数
	responseWriter := httptest.NewRecorder()

	// 模拟请求的返回结果
	helloHandler(responseWriter, req)

	// 获得响应结果
	bytes, _ := io.ReadAll(responseWriter.Result().Body)

	if string(bytes) != "hello world" {
		t.Fatal("expected hello world, but got", string(bytes))
	}
}
