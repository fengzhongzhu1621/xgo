package nethttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return req
}

// RequestWithBody 创建一个新的HTTP请求
func RequestWithBody2(method, path, body string) (req *http.Request) {
	// http.NewRequest 创建一个新的HTTP请求
	// method：HTTP请求方法，如"GET"、"POST"、"PUT"、"DELETE"等。
	// url：请求的目标URL。
	// body：请求的主体内容，通常用于POST、PUT等需要发送数据的请求。如果请求不需要主体内容，可以传入nil或ioutil.NopCloser(io.Discard)

	// 将大字符串转换为一个io.Reader
	// strings.NewReader将一个大字符串转换成一个实现了io.Reader接口的对象。
	// 这个io.Reader对象可以逐步读取数据并写入到请求的body中。
	reader := strings.NewReader(body)
	req, _ = http.NewRequest(method, path, reader)

	return req
}

func RequestWithFile(method, path, filepath string) (req *http.Request) {
	// 打开一个大文件
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	req, _ = http.NewRequest(method, path, file)

	return req
}

// ReadRequestBody will return the body in []byte, without change the origin body
func ReadRequestBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, ErrNilRequestBody
	}

	// 将请求体转换为 []byte
	byt, err := io.ReadAll(r.Body)
	// io.NopCloser 是 Go 语言标准库 io 包中的一个函数，它接受一个 io.Reader 接口类型的参数，并返回一个新的 io.ReadCloser 接口类型的值。
	// 这个新的 io.ReadCloser 实现了 Close 方法，但该方法不执行任何操作（即“nop”表示无操作）。
	// io.NopCloser 的主要用途是将一个只读的 io.Reader 转换为一个 io.ReadCloser，以便在需要 io.ReadCloser 的地方使用，而不需要实际关闭底层资源。
	r.Body = io.NopCloser(bytes.NewReader(byt))

	return byt, err
}

// PeekRequest “窥视”（Peek）一个 HTTP 请求的主体内容，
// 即在不消耗请求主体的情况下读取其内容。
// 这在需要在处理请求之前查看请求体内容的场景中非常有用，例如日志记录、验证或其他预处理步骤。
func PeekRequest(req *http.Request) ([]byte, error) {
	if req.Body != nil {
		// 读取 req.Body 中的所有数据。io.ReadAll 会一直读取直到遇到 EOF（文件结束符）或发生错误
		byt, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		// 读取完请求主体后，原始的 req.Body 已经被消耗（即处于 EOF 状态）。
		// 为了不影响后续对请求主体的正常处理（例如，路由处理器仍然需要读取请求体），需要将 req.Body 恢复到一个可读取的状态。
		//
		// 这对于 http.Request.Body 是必要的，因为 Body 需要是一个 ReadCloser。
		req.Body = io.NopCloser(bytes.NewBuffer(byt))
		return byt, nil
	}
	return make([]byte, 0), nil
}
