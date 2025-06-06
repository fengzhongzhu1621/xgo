package cast

import (
	"reflect"
	"unsafe"
)

// //////////////////////////////////////////////////////////////////////////////////////

// StringToBytes 字符串转换为[]bytes converts string to byte slice without a memory allocation.
// StringHeader 和 SliceHeader 在 Go 1.20+ 被标记为弃用
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

// SafeStringToBytes 字符串转换为字节数组，标准库方法，安全且高效，推荐
func SafeStringToBytes(s string) []byte {
	return []byte(s)
}

func SafeStringCopyToBytes(s string) []byte {
	b := make([]byte, len(s))
	copy(b, s)
	return b
}
