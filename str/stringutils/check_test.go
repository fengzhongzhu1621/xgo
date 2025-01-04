package stringutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestIsString 检查值的数据类型是否为字符串。
func TestIsString(t *testing.T) {
	result1 := strutil.IsString("")
	result2 := strutil.IsString("a")
	result3 := strutil.IsString(1)
	result4 := strutil.IsString(true)
	result5 := strutil.IsString([]string{"a"})

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, false, result3)
	assert.Equal(t, false, result4)
	assert.Equal(t, false, result5)
}

// IsBlank 检查一个字符串是空白还是空。
// func IsBlank(str string) bool
func TestIsBlank(t *testing.T) {
	result1 := strutil.IsBlank("")
	result2 := strutil.IsBlank("\t\v\f\n")
	result3 := strutil.IsBlank(" 中文")

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, false, result3)
}

// IsNotBlank 检查一个字符串是否不是空白或者不为空。
// func IsNotBlank(str string) bool
func TestIsNotBlank(t *testing.T) {
	result1 := strutil.IsNotBlank("")
	result2 := strutil.IsNotBlank("    ")
	result3 := strutil.IsNotBlank("\t\v\f\n")
	result4 := strutil.IsNotBlank(" 中文")
	result5 := strutil.IsNotBlank("    world    ")

	assert.Equal(t, false, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, false, result3)
	assert.Equal(t, true, result4)
	assert.Equal(t, true, result5)
}

// HasPrefixAny 检查一个字符串是否以一组指定的字符串中的任意一个开头。
// func HasPrefixAny(str string, prefixes []string) bool
func TestHasPrefixAny(t *testing.T) {
	result1 := strutil.HasPrefixAny("foo bar", []string{"fo", "xyz", "hello"})
	result2 := strutil.HasPrefixAny("foo bar", []string{"oom", "world"})

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}

// HasSuffixAny 检查一个字符串是否以一组指定的字符串中的任意一个结尾。
// func HasSuffixAny(str string, prefixes []string) bool
func TestHasSuffixAny(t *testing.T) {
	result1 := strutil.HasSuffixAny("foo bar", []string{"bar", "xyz", "hello"})
	result2 := strutil.HasSuffixAny("foo bar", []string{"oom", "world"})

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}
