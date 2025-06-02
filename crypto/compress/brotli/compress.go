package brotli

import (
	"bytes"

	"github.com/google/brotli/go/cbrotli"
)

// Compress 使用Brotli算法压缩数据
// 参数:
//
//	data - 要压缩的字节数据
//
// 返回:
//
//	压缩后的字节数据
//	可能的错误
func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := cbrotli.NewWriter(&buf, cbrotli.WriterOptions{Quality: 11})
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decompress 使用Brotli算法解压缩数据
// 参数:
//
//	data - 压缩后的字节数据
//
// 返回:
//
//	解压后的字节数据
//	可能的错误
func Decompress(data []byte) ([]byte, error) {
	reader := cbrotli.NewReader(bytes.NewReader(data))
	var buf bytes.Buffer
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 示例用法
func main() {
	original := []byte("This is some sample text to compress and decompress using Brotli.")

	// 压缩
	compressed, err := Compress(original)
	if err != nil {
		panic(err)
	}

	// 解压缩
	decompressed, err := Decompress(compressed)
	if err != nil {
		panic(err)
	}

	// 验证
	if string(original) != string(decompressed) {
		panic("Decompressed data doesn't match original")
	}
}
