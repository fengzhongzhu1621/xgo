package stringutils

import (
	"strconv"
)

func IsPrint(s string) bool {
	for _, c := range s {
		if !strconv.IsPrint(c) {
			return false
		}
	}

	return true
}

// QuoteIfNeeded 检查给定的字符串 s 是否包含可打印字符。
// 如果字符串中包含不可打印的字符（例如控制字符），则使用 strconv.Quote 将其转义并加上引号；
// 否则，直接返回原始字符串。
func QuoteIfNeeded(s string) string {
	if !IsPrint(s) {
		return strconv.Quote(s)
	}

	return s
}

func QuoteIfNeededV(s []string) []string {
	ret := make([]string, len(s))

	for i, v := range s {
		ret[i] = QuoteIfNeeded(v)
	}

	return ret
}

// 将普通字符串转换为带有转义字符的带引号的字符串字面值
func QuoteV(s []string) []string {
	ret := make([]string, len(s))

	for i, v := range s {
		ret[i] = strconv.Quote(v)
	}

	return ret
}

// UnquoteIfPossible 用于将带有转义字符的带引号的字符串字面值转换为普通字符串
func UnquoteIfPossible(s string) (string, error) {
	if len(s) == 0 || s[0] != '"' {
		return s, nil
	}

	return strconv.Unquote(s)
}
