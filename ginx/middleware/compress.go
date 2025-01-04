package middleware

// https://github.com/cheikhsimsol/7zip

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type CompressFunc func(data []byte) ([]byte, error)

// CompressResponseWriter is a custom writer that
// mimics Gin's ResponseWriter without actual I/O.
type CompressResponseWriter struct {
	headers            http.Header
	buffer             *bytes.Buffer
	status             int
	gin.ResponseWriter // risky as this value is nil and
	// will trigger an error if an um-implemented method is called.
}

func (w *CompressResponseWriter) Header() http.Header {
	return w.headers
}

func (w *CompressResponseWriter) Write(data []byte) (int, error) {
	return w.buffer.Write(data)
}

func (w *CompressResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func (w *CompressResponseWriter) Status() int {
	return w.status
}

func (w *CompressResponseWriter) Written() bool {
	return w.buffer.Len() > 0
}

// Compress7zip 使用 7z 二进制程序压缩内容
func Compress7zip(data []byte) ([]byte, error) {
	// Create a temporary file to store the input data.
	tmpFile, err := os.CreateTemp("", "input_*.tmp")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name()) // Ensure the temp file is deleted after use.

	// Write the data to the temporary file.
	_, err = tmpFile.Write(data)
	if err != nil {
		return nil, err
	}
	tmpFile.Close()

	// Create a temporary directory for the compressed file output.
	tmpDir := filepath.Join(os.TempDir(), "compress_output")
	err = os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir) // Cleanup the output directory.

	// Output file path.
	outputFile := filepath.Join(tmpDir, "output.7z")

	// Run the 7z command to compress the file.
	cmd := exec.Command("7z", "a", "-mx=9", outputFile, tmpFile.Name())

	err = cmd.Run()
	if err != nil {
		return nil, errors.New("failed to run 7z command: " + err.Error())
	}

	// Read the compressed file into memory.
	compressedData, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, err
	}

	return compressedData, nil
}

func CompressResponse(cf CompressFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		buffer := bytes.Buffer{}
		compressWriter := &CompressResponseWriter{
			headers: http.Header{},
			buffer:  &buffer,
			status:  http.StatusOK,
		}

		originalWriter := c.Writer
		c.Writer = compressWriter

		// will return after all
		// subsequent middleware/handlers have executed.
		// 在所有后续的中间件和处理器执行完毕后返回
		c.Next()

		// io.Copy(c.Writer, file) 将结果写入到响应
		// 此时响应内容已经调用 CompressResponseWriter.Write() 方法写入到 buffer 中

		// 压缩文件
		result := buffer.Bytes()
		compressed, err := cf(result)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// 压缩完成后，我会重置原始的ResponseWriter
		c.Writer = originalWriter

		// 告知客户端这是一个压缩文件
		if compressWriter.status == http.StatusOK {
			c.Writer.Header().Set("Content-Type", "application/x-7z-compressed")
		}

		// 将压缩后的字节写入响应
		c.Writer.WriteHeader(compressWriter.status)
		c.Writer.Write(compressed)
	}
}
