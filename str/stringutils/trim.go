package stringutils

import "strings"

// TrimRight 去掉字符串的后缀.
func TrimRight(str string, substring string) string {
	idx := strings.LastIndex(str, substring)
	if idx < 0 {
		return str
	}
	return str[:idx]
}

// TrimLeft 去掉字符串的前缀.
func TrimLeft(str string, substring string) string {
	return strings.TrimPrefix(str, substring)
}
