package stringutils

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
)

// TestFieldsFunc 测试自定义分割
func TestFieldsFunc(t *testing.T) {
	s := "Hello, 世界! This is a test."

	// 用于将字符串按照指定的分隔函数进行分割，并返回一个字符串切片。这个函数非常灵活，因为它允许你自定义分隔符的判断逻辑。
	// 使用 FieldsFunc 分割字符串，以空白字符（空格、制表符等）和标点符号作为分隔符
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})

	fmt.Println(fields) // 输出：["Hello" "世界" "This" "is" "a" "test"]
}
