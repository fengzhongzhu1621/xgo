package stringutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gookit/goutil/arrutil"
	"github.com/stretchr/testify/assert"
)

// TestTrim 从字符串的开头和结尾去除空格（或其他字符）。可选参数 characterMask 指定要去除的附加字符。
// 先去掉附加字符，然后再去掉空白字符
// func Trim(str string, characterMask ...string) string
func TestTrim(t *testing.T) {
	is := assert.New(t)

	{
		result1 := strutil.Trim("\nabcd")

		str := "$ ab    cd $ "

		result2 := strutil.Trim(str)
		result3 := strutil.Trim(str, "$")
		result4 := strutil.Trim("     hi,       你好 ", "你好")

		assert.Equal(t, "abcd", result1)
		assert.Equal(t, "$ ab    cd $", result2)
		assert.Equal(t, "ab    cd", result3)
		assert.Equal(t, "hi,", result4)
	}

	{
		ss := arrutil.TrimStrings([]string{" a", "b ", " c "})
		is.Equal("[a b c]", fmt.Sprint(ss))
		ss = arrutil.TrimStrings([]string{",a", "b.", ",.c,"}, ",.")
		is.Equal("[a b c]", fmt.Sprint(ss))
		ss = arrutil.TrimStrings([]string{",a", "b.", ",.c,"}, ",", ".")
		is.Equal("[a b c]", fmt.Sprint(ss))
	}
}

// TestSplitAndTrim 通过字符串delimiter将字符串str分割为切片，并对切片的每个元素调用Trim。它会忽略在Trim后为空的元素。
// 先去掉附加字符，然后再去掉空白字符
// func SplitAndTrim(str, delimiter string, characterMask ...string) []string
func TestSplitAndTrim(t *testing.T) {
	str := "  a,b, c,d, $1  "

	result1 := strutil.SplitAndTrim(str, ",")
	result2 := strutil.SplitAndTrim(str, ",", "$")
	result3 := strutil.SplitAndTrim("     hi,       你好 ", "你好")

	assert.Equal(t, []string{"a", "b", "c", "d", "$1"}, result1)
	assert.Equal(t, []string{"a", "b", "c", "d", "1"}, result2)
	assert.Equal(t, []string{"hi,"}, result3)
}
