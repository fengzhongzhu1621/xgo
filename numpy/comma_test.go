package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/formatter"
)

// TestComma 从右向左每3个数字为一个单位，在数字值中添加逗号，并在前面加上一个前缀符号字符。如果值是像“aa”这样的无效数字字符串，则返回空字符串。
// func Comma[T constraints.Float | constraints.Integer | string](value T, prefixSymbol string) string
func TestComma(t *testing.T) {
	result1 := formatter.Comma("123", "")
	result2 := formatter.Comma("12345", "$")
	result3 := formatter.Comma(1234567, "￥")

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 123
	// $12,345
	// ￥1,234,567
}
