package req

import (
	"testing"

	"github.com/imroc/req/v3"
)

func TestEnableForceHTTP1(t *testing.T) {
	req.DevMode() // Treat the package name as a Client, enable development mode
	req.MustGet(
		"https://httpbin.org/uuid",
	) // Treat the package name as a Request, send GET request.

	// HTTP/1.1 相对于 HTTP/1 新增了许多请求方法，现今我们常用的PUT、PATCH、DELETE、CONNECT、TRACE 和OPTIONS 等都是在 HTTP/1.1 时新增的。
	req.EnableForceHTTP1() // Force using HTTP/1.1

	req.MustGet("https://httpbin.org/uuid")
}
