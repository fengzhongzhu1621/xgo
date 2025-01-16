package cast

import (
	"reflect"
	"unsafe"
)

// //////////////////////////////////////////////////////////////////////////////////////

// StringToBytes 字符串转换为[]bytes converts string to byte slice without a memory allocation.
// 效率更高.
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// Bytes converts stringutils to byte slice.
func Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// rawStrToBytes 字符串转换为字节数组
func rawStrToBytes(s string) []byte {
	return []byte(s)
}

// SafeBytes 字符串转换为字节数组
func SafeBytes(s string) []byte {
	return []byte(s)
}

// TruncateBytes 截断字节切片
func TruncateBytes(content []byte, length int) []byte {
	if len(content) > length {
		return content[:length]
	}
	return content
}
