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

func SublimeContains(s, substr string) bool {
	// 将字符串转换为 rune 切片，您可以逐个访问和处理字符串中的每个字符，而无需担心底层编码细节
	rs, rsubstr := []rune(s), []rune(substr)
	// 判断子串的长度
	if len(rsubstr) > len(rs) {
		return false
	}

	var ok = true
	var i, j = 0, 0
	for ; i < len(rsubstr); i++ {
		found := -1
		for ; j < len(rs); j++ {
			if rsubstr[i] == rs[j] {
				found = j
				break
			}
		}
		if found == -1 {
			ok = false
			break
		}
		j += 1
	}
	return ok
}