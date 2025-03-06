package nethttp

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

var (
	ErrNilRequestBody = errors.New("request Body is nil")
)

// RequestWithBody 创建一个新的HTTP请求
func RequestWithBody(method, path, body string) (req *http.Request) {
	// http.NewRequest 创建一个新的HTTP请求
	// method：HTTP请求方法，如"GET"、"POST"、"PUT"、"DELETE"等。
	// url：请求的目标URL。
	// body：请求的主体内容，通常用于POST、PUT等需要发送数据的请求。如果请求不需要主体内容，可以传入nil或ioutil.NopCloser(io.Discard)
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body)) //nolint
	return req
}

// ReadRequestBody will return the body in []byte, without change the origin body
func ReadRequestBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, ErrNilRequestBody
	}

	// 将请求体转换为 []byte
	body, err := io.ReadAll(r.Body)
	// io.NopCloser 是 Go 语言标准库 io 包中的一个函数，它接受一个 io.Reader 接口类型的参数，并返回一个新的 io.ReadCloser 接口类型的值。
	// 这个新的 io.ReadCloser 实现了 Close 方法，但该方法不执行任何操作（即“nop”表示无操作）。
	// io.NopCloser 的主要用途是将一个只读的 io.Reader 转换为一个 io.ReadCloser，以便在需要 io.ReadCloser 的地方使用，而不需要实际关闭底层资源。
	r.Body = io.NopCloser(bytes.NewReader(body))

	return body, err
}
