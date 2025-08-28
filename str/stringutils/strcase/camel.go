package strcase

import (
	"strings"
)

// Converts a string to CamelCase.
func toCamelInitCase(s string, initCase bool) string {
	// 去掉空白字符
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	// 自定义缩写转换
	if a, ok := uppercaseAcronym[s]; ok {
		s = a
	}
	// 预先给字符串分配内存，实现分配足够的内容，减少运行时的内存分配
	n := strings.Builder{}
	n.Grow(len(s))

	capNext := initCase
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}

// ToCamel converts a string to CamelCase.
func ToCamel(s string) string {
	return toCamelInitCase(s, true)
}

// ToLowerCamel converts a string to lowerCamelCase.
func ToLowerCamel(s string) string {
	return toCamelInitCase(s, false)
}
