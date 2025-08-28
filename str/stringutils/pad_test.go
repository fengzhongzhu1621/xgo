package stringutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestPad 如果字符串的长度小于指定的尺寸，就在其左右两侧进行填充。
func TestPad(t *testing.T) {
	result1 := strutil.Pad("foo", 1, "bar")
	result2 := strutil.Pad("foo", 2, "bar")
	result3 := strutil.Pad("foo", 3, "bar")
	result4 := strutil.Pad("foo", 4, "bar")
	result5 := strutil.Pad("foo", 5, "bar")
	result6 := strutil.Pad("foo", 6, "bar")
	result7 := strutil.Pad("foo", 7, "bar")
	result8 := strutil.Pad("foo", 8, "bar")
	result9 := strutil.Pad("foo", 9, "bar")
	result10 := strutil.Pad("foo", 7, "*")

	assert.Equal(t, "foo", result1)
	assert.Equal(t, "foo", result2)
	assert.Equal(t, "foo", result3)
	assert.Equal(t, "foob", result4)
	assert.Equal(t, "bfoob", result5)
	assert.Equal(t, "bfooba", result6)
	assert.Equal(t, "bafooba", result7)
	assert.Equal(t, "bafoobar", result8)
	assert.Equal(t, "barfoobar", result9)
	assert.Equal(t, "**foo**", result10)
}

// TestPadEnd 如果字符串的长度小于指定的尺寸，就在右侧填充字符串。
func TestPadEnd(t *testing.T) {
	result1 := strutil.PadEnd("foo", 1, "bar")
	result2 := strutil.PadEnd("foo", 2, "bar")
	result3 := strutil.PadEnd("foo", 3, "bar")
	result4 := strutil.PadEnd("foo", 4, "bar")
	result5 := strutil.PadEnd("foo", 5, "bar")
	result6 := strutil.PadEnd("foo", 6, "bar")
	result7 := strutil.PadEnd("foo", 7, "bar")
	result8 := strutil.PadEnd("foo", 8, "bar")
	result9 := strutil.PadEnd("foo", 9, "bar")

	assert.Equal(t, "foo", result1)
	assert.Equal(t, "foo", result2)
	assert.Equal(t, "foo", result3)
	assert.Equal(t, "foob", result4)
	assert.Equal(t, "fooba", result5)
	assert.Equal(t, "foobar", result6)
	assert.Equal(t, "foobarb", result7)
	assert.Equal(t, "foobarba", result8)
	assert.Equal(t, "foobarbar", result9)
}

// TestPadStart 如果字符串的长度小于指定的尺寸，就在其左侧填充。
func TestPadStart(t *testing.T) {
	result1 := strutil.PadStart("foo", 1, "bar")
	result2 := strutil.PadStart("foo", 2, "bar")
	result3 := strutil.PadStart("foo", 3, "bar")
	result4 := strutil.PadStart("foo", 4, "bar")
	result5 := strutil.PadStart("foo", 5, "bar")
	result6 := strutil.PadStart("foo", 6, "bar")
	result7 := strutil.PadStart("foo", 7, "bar")
	result8 := strutil.PadStart("foo", 8, "bar")
	result9 := strutil.PadStart("foo", 9, "bar")

	assert.Equal(t, "foo", result1)
	assert.Equal(t, "foo", result2)
	assert.Equal(t, "foo", result3)
	assert.Equal(t, "bfoo", result4)
	assert.Equal(t, "bafoo", result5)
	assert.Equal(t, "barfoo", result6)
	assert.Equal(t, "barbfoo", result7)
	assert.Equal(t, "barbafoo", result8)
	assert.Equal(t, "barbarfoo", result9)
}

// TestHideString 使用参数 replaceChar 隐藏源字符串中的某些字符。替换范围是 origin[start : end]。[start, end)。
// func HideString(origin string, start, end int, replaceChar string) string
func TestHideString(t *testing.T) {
	str := "13242658976"

	result1 := strutil.HideString(str, 3, 3, "*")
	result2 := strutil.HideString(str, 3, 4, "*")
	result3 := strutil.HideString(str, 3, 7, "*")
	result4 := strutil.HideString(str, 7, 11, "*")
	result5 := strutil.HideString(str, 7, 100, "*")

	assert.Equal(t, "13242658976", result1)
	assert.Equal(t, "132*2658976", result2)
	assert.Equal(t, "132****8976", result3)
	assert.Equal(t, "1324265****", result4)
	assert.Equal(t, "1324265****", result5)
}
