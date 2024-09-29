package stringutils

import (
	"reflect"
	"strings"
)

// Deprecated: 翻转切片 panic if s is not a slice.
func ReflectReverseSlice(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// Deprecated: 翻转切片，返回一个新的切片，有copy的耗损.
func ReverseSliceGetNew(s []string) []string {
	a := make([]string, len(s))
	copy(a, s)

	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}

	return a
}

// ReverseSlice 翻转切片，值会改变，性能最高.
func ReverseSlice(a []string) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

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
