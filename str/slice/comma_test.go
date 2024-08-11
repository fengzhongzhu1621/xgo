package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/formatter"
	"github.com/stretchr/testify/assert"
)

// 用逗号每隔3位分割数字/字符串，支持前缀添加符号。参数 value 必须是数字或者可以转为数字的字符串, 否则返回空字符串。
func TestComma(t *testing.T) {
	result1 := formatter.Comma("123", "")
	result2 := formatter.Comma("12345", "$")
	result3 := formatter.Comma(1234567, "￥")

	assert.Equal(t, result1, "123")
	assert.Equal(t, result2, "$12,345")
	assert.Equal(t, result3, "￥1,234,567")
}
