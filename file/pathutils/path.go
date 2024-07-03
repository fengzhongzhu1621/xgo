package pathutils

import (
	"os"
)

type PathInfo struct {
	Name  string // 路径名称
	IsDir bool   // 是否是目录
}


// 获得最后一个字符.
func LastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

// GetWd 获得应用程序当前路径.
func GetWd() string {
	wd, _ := os.Getwd()
	return wd
}

// bufApp 将字符串 s[:w]拷贝到缓存 buf 中，重新计算缓存的大小刚好能容纳字符串的长度
// 并修改指定索引位置的字符值
//
// buf 为指针的原因是这是一个可修改值，执行后会修改字节数组的内容
func bufApp(buf *[]byte, s string, w int, c byte) {
	b := *buf
	if len(b) == 0 {
		// No modification of the original string so far.
		// If the next character is the same as in the original string, we do
		// not yet have to allocate a buffer.
		if s[w] == c {
			return
		}

		// 重新计算缓存的大小
		// Otherwise use either the stack buffer, if it is large enough, or
		// allocate a new buffer on the heap, and copy all previous characters.
		length := len(s)
		if length > cap(b) {
			// 源字符串超过缓存 bug 的长度，则重新创建一个新的数组
			*buf = make([]byte, length)
		} else {
			// 否则减少缓存的大小
			*buf = (*buf)[:length]
		}
		b = *buf

		// 将字符串 [ 0, w) 拷贝到缓存中
		copy(b, s[:w])
	}

	// 并将 w 索引修改为指定字符 c
	b[w] = c
}
