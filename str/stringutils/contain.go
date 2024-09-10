package stringutils

import (
	"sort"
	"strings"
)

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
	// 先排序再搜索
	sort.Strings(src)
	index := sort.SearchStrings(src, target)
	if index < len(src) && src[index] == target {
		return true
	}
	return false
}

// StringsMatch 如果 key 是空字符串，则认为它匹配任何内容
// 如果 key 是空字符串，则认为它匹配任何内容。
func StringsMatch(key string, matchList ...string) bool {
	if key == "" {
		return true
	}
	return strings.Contains(strings.Join(matchList, "#"), key)
}

// StringsMatchObscure 检查 key 是否包含在由 matchList 中的字符串拼接而成的字符串中，且在比较时不区分大小写。
// 如果 key 是空字符串，则认为它匹配任何内容
func StringsMatchObscure(key string, matchList ...string) bool {
	var (
		t string
	)
	if key == "" {
		return true
	}

	for _, v := range matchList {
		t += strings.ToLower(v)
	}
	return strings.Contains(t, strings.ToLower(key))
}

// SublimeContains 判断是否包含子串
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
