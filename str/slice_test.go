package str

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/formatter"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// 判断 slice 是否包含指定值
func TestContain(t *testing.T) {
	result1 := slice.Contain([]string{"a", "b", "c"}, "a")
	result2 := slice.Contain([]int{1, 2, 3}, 4)

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}

// 按照 size 参数均分 slice
func TestChunk(t *testing.T) {
	arr := []string{"a", "b", "c", "d", "e"}

	result1 := slice.Chunk(arr, 1)
	result2 := slice.Chunk(arr, 2)

	fmt.Println(result1)
	fmt.Println(result2)
}

// 检查两个切片是否相等，相等条件：切片长度相同，元素顺序和值都相同。
func TestEqual(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{1, 3, 2}

	result1 := slice.Equal(s1, s2)
	result2 := slice.Equal(s1, s3)

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}

// 用逗号每隔3位分割数字/字符串，支持前缀添加符号。参数 value 必须是数字或者可以转为数字的字符串, 否则返回空字符串。
func TestComma(t *testing.T) {
	result1 := formatter.Comma("123", "")
	result2 := formatter.Comma("12345", "$")
	result3 := formatter.Comma(1234567, "￥")

	assert.Equal(t, result1, "123")
	assert.Equal(t, result2, "$12,345")
	assert.Equal(t, result3, "￥1,234,567")
}
