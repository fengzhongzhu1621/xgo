package string_utils

import "strings"

/**
 * 根据分隔符分割字符串
 */
func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}
