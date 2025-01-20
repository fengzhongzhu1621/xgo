package stringutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestReplaceWithMap 返回str的副本，该副本以无序的方式、区分大小写地被映射替换
// func ReplaceWithMap(str string, replaces map[string]string) string
func TestReplaceWithMap(t *testing.T) {
	str := "ac ab ab ac"
	replaces := map[string]string{
		"a": "1",
		"b": "2",
	}

	result := strutil.ReplaceWithMap(str, replaces)

	assert.Equal(t, "1c 12 12 1c", result)
}

// TestRemoveWhiteSpace 从字符串中删除空白字符。当设置replaceAll为true时，删除所有空白，false则仅将连续的空白字符替换为一个空格。
// func RemoveWhiteSpace(str string, repalceAll bool) string
func TestRemoveWhiteSpace(t *testing.T) {
	str := " hello   \r\n    \t   world"

	result1 := strutil.RemoveWhiteSpace(str, true)
	result2 := strutil.RemoveWhiteSpace(str, false)

	assert.Equal(t, "helloworld", result1)
	assert.Equal(t, "hello world", result2)
}

func TestDot(t *testing.T) {
	str := "test1.test2.test3.test4"
	encodedStr := EncodeDot(str)
	assert.Equal(t, encodedStr, "test1\\u002etest2\\u002etest3\\u002etest4")

	decodedStr := DecodeDot(encodedStr)
	assert.Equal(t, decodedStr, str)
}
