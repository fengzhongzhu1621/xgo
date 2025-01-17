package stringutils

import (
	"strings"
)

// ReverseString 字符串反转
func ReverseString(s string) string {
	if s == "" {
		return ""
	}

	var newString []string
	for i := len(s) - 1; i >= 0; i-- {
		newString = append(newString, string(s[i]))
	}
	return strings.Join(newString, "")
}
