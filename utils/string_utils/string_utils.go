package string_utils

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

/**
 * 根据分隔符分割字符串
 */
func Head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}

////////////////////////////////////////////////////////
// 字符串拼接

// Deprecated: 低效的字符串拼接
func StringPlus(p []string) string {
	var s string
	l := len(p)
	for i := 0; i < l; i++ {
		s += p[i]
	}
	return s
}

// Deprecated: 低效的字符串拼接
func StringFmt(p []interface{}) string {
	return fmt.Sprint(p...)
}

// Deprecated: 低效的字符串拼接
func StringBuffer(p []string) string {
	var b bytes.Buffer
	l := len(p)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

func StringJoin(p []string) string {
	return strings.Join(p, "")
}

func StringBuilder(p []string) string {
	var b strings.Builder
	l := len(p)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

func StringBuilderEx(p []string, cap int) string {
	var b strings.Builder
	l := len(p)
	// 实现分配足够的内容，减少运行时的内存分配
	b.Grow(cap)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

/**
 * 切片去重
 * 空struct不占内存空间，使用它来实现我们的函数空间复杂度是最低的。
 */
func RemoveDuplicateElement(items []string) []string {
	result := make([]string, 0, len(items))
	// 定义集合
	set := map[string]struct{}{}
	for _, item := range items {
		if _, ok := set[item]; !ok {
			// 如何集合中不存在此元素，则加入集合
			set[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

/**
 * 去掉字符串的后缀
 */
func TrimRight(str string, substring string) string {
	idx := strings.LastIndex(str, substring)
	if idx < 0 {
		return str
	}
	return str[:idx]
}

/**
 * 去掉字符串的前缀
 */
func TrimLeft(str string, substring string) string {
	return strings.TrimPrefix(str, substring)
}

func UnicodeTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToUpper(v)) + s[k+1:]
	}
	return ""
}

func UnicodeUnTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToLower(v)) + s[k+1:]
	}
	return ""
}

/**
 * 返回数组最后一个元素
 */
func Last(list []string) string {
	return list[len(list) - 1]
}
