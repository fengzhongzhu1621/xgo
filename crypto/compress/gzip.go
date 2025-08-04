package compress

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Gzip 将输入的字节切片压缩为 gzip 格式，并返回压缩后的字节切片
func Gzip(in []byte) ([]byte, error) {
	if in == nil {
		return nil, nil
	}

	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)

	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}

// UnGZip 将输入的 gzip 格式的字节切片解压缩，并返回解压缩后的字节切片
func UnGZip(in []byte) ([]byte, error) {
	if in == nil {
		return nil, nil
	}

	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// IsGZIP 判断输入的字股切片是否为 gzip 根据格式
func IsGZIP(in []byte) bool {
	if in == nil {
		return false
	}

	// gzip 文件头至少有 10 个字节
	if l := len(in); l < 10 {
		return false
	}

	// 比较输入的字节切片的前两个字节与 gzip 文件头的标识（0x1F, 0x8B）。
	// 如果相等，说明输入的字股切片是 gzip 根据格式
	if bytes.Equal(in[:2], []byte{0x1F, 0x8B}) {
		return true
	}
	return false
}
