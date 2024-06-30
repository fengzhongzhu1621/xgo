package stringutils

import (
	"unicode/utf8"
)

// ChangeInitialCase 将字符串的首字符转换为指定格式
func ChangeInitialCase(s string, mapper func(rune) rune) string {
	if s == "" {
		return s
	}
	// 返回第一个utf8字符，n是字符的长度，即返回首字符
	r, n := utf8.DecodeRuneInString(s)
	// 根据mapper方法转换首字符
	return string(mapper(r)) + s[n:]
}
