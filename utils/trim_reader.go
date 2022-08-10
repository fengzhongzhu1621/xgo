package utils

import (
	"bytes"
	"io"
	"unicode"
)

// trimReader is a custom io.Reader that will trim any leading
// whitespace, as this can cause email imports to fail.
type TrimReader struct {
	Rd      io.Reader
	Trimmed bool // 标记是否存在空白字符被去掉
}

// Read trims off any unicode whitespace from the originating reader
func (tr *TrimReader) Read(buf []byte) (int, error) {
	// 从输入读取内容保存到buf中，返回读取的字节数
	n, err := tr.Rd.Read(buf)
	if err != nil {
		return n, err
	}
	if !tr.Trimmed {
		// 去掉左侧空白字符
		t := bytes.TrimLeftFunc(buf[:n], unicode.IsSpace)
		tr.Trimmed = true
		n = copy(buf, t)
	}
	return n, err
}
