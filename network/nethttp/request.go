package nethttp

import (
	"bytes"
	"net/http"
)

func RequestWithBody(method, path, body string) (req *http.Request) {
	// http.NewRequest 创建一个新的HTTP请求
	// method：HTTP请求方法，如"GET"、"POST"、"PUT"、"DELETE"等。
	// url：请求的目标URL。
	// body：请求的主体内容，通常用于POST、PUT等需要发送数据的请求。如果请求不需要主体内容，可以传入nil或ioutil.NopCloser(io.Discard)
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body)) // nolint
	return req
}
