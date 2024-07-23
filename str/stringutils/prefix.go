package stringutils

import (
	"regexp"
	"strings"
)

func FixPrefix(prefix string) string {
	// 替换所有匹配的子字符串，在 prefix 中去掉 /*$，例如 abc/d -> abc
	prefix = regexp.MustCompile(`/*$`).ReplaceAllString(prefix, "")
	// abc -> /abc
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	if prefix == "/" {
		prefix = ""
	}
	return prefix
}
