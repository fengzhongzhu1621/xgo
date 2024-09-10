package stringutils

import (
	"unicode"
	"unicode/utf8"

	"github.com/fengzhongzhu1621/xgo/str/bytesconv"
)

// ToLower 字符串转换为小写，在转化前先判断是否包含大写字符，比strings.ToLower性能高.
func ToLower(s string) string {
	// 判断字符串是否包含小写字母
	if IsLower(s) {
		return s
	}
	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	// []bytes转换为字符串
	return bytesconv.BytesToString(b)
}

// UnicodeTitle 首字母大写.
func UnicodeTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToUpper(v)) + s[k+1:]
	}
	return ""
}

// UnicodeUnTitle 首字母小写.
func UnicodeUnTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToLower(v)) + s[k+1:]
	}
	return ""
}

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
