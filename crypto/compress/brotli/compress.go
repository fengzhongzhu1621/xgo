//go:build brotli
// +build brotli

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
