package stringutils

import (
	"sort"
)

// IsLower 判断字符串是否包含小写字母.
func IsLower(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

// StringInSlice 判断字符串是否在切片中.
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// In 判断字符串是否在数组中
func In(target string, src []string) bool {
	sort.Strings(src)
	index := sort.SearchStrings(src, target)
	if index < len(src) && src[index] == target {
		return true
	}
	return false
}

