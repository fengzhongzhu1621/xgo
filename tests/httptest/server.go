package httptest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// NewHTTPServer 创建一个 HTTP 服务器，用于测试。
func NewHTTPServer(resp string) *httptest.Server {
	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, resp)
	}))
}
