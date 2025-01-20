package stringutils

import (
	"strings"
)

// Head 根据分隔符分割字符串.
func Head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}

func SplitString(r rune) bool {
	return r == ';' || r == ',' || r == '\n'
}

// SplitMulti split multi string by sep string.
func SplitMulti(ss []string, sep string) []string {
	ns := make([]string, 0, len(ss)+1)
	for _, s := range ss {
		ns = append(ns, strings.Split(s, sep)...)
	}
	return ns
}
