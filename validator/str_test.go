package validator

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/duke-git/lancet/v2/validator"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/simplifiedchinese"
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

// TestIsEmptyString 检查字符串是否为空。
// func IsEmptyString(str string) bool
func TestIsEmptyString(t *testing.T) {
	result1 := validator.IsEmptyString("")
	result2 := validator.IsEmptyString(" ")
	result3 := validator.IsEmptyString("\t")

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, false, result3)
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

// TestContainChinese 验检查字符串是否包含普通话中文。
func TestContainChinese(t *testing.T) {
	result1 := validator.ContainChinese("你好")
	result2 := validator.ContainChinese("你好 hello")
	result3 := validator.ContainChinese("hello")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, true)
	assert.Equal(t, result3, false)
}

// TestContainLetter 检查该字符串是否至少包含一个字母。
func TestContainLetter(t *testing.T) {
	result1 := validator.ContainLetter("你好")
	result2 := validator.ContainLetter("&@#$%^&*")
	result3 := validator.ContainLetter("ab1")

	assert.Equal(t, result1, false)
	assert.Equal(t, result2, false)
	assert.Equal(t, result3, true)
}

// TestContainLower 检查字符串是否至少包含一个小写字母a-z。
func TestContainLower(t *testing.T) {
	result1 := validator.ContainLower("abc")
	result2 := validator.ContainLower("aBC")
	result3 := validator.ContainLower("ABC")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, true)
	assert.Equal(t, result3, false)
}

// TestContainUpper 检查字符串是否至少包含一个大写字母 A-Z。
func TestContainUpper(t *testing.T) {
	result1 := validator.ContainUpper("ABC")
	result2 := validator.ContainUpper("abC")
	result3 := validator.ContainUpper("abc")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, true)
	assert.Equal(t, result3, false)
}

// TestIsChineseMobile 验证字符串是否是中国手机号码
func TestIsChineseMobile(t *testing.T) {
	result1 := validator.IsChineseMobile("17530367777")
	result2 := validator.IsChineseMobile("121212121")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}

// TestIsChinesePhone 验证字符串是否是中国电话座机号码
func TestIsChinesePhone(t *testing.T) {
	result1 := validator.IsChinesePhone("010-32116675")
	result2 := validator.IsChinesePhone("123-87562")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}

// TestIsAlpha 检查字符串是否仅包含字母（a-zA-Z）。
func TestIsAlpha(t *testing.T) {
	result1 := validator.IsAlpha("abc")
	result2 := validator.IsAlpha("ab1")
	result3 := validator.IsAlpha("")

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, false, result2, "result2")
	assert.Equal(t, false, result3, "result3")
}

// TestIsAllUpper 检查字符串是否全部由大写字母 A-Z 组成。
// func IsAllUpper(str string) bool
func TestIsAllUpper(t *testing.T) {
	result1 := validator.IsAllUpper("ABC")
	result2 := validator.IsAllUpper("ABc")
	result3 := validator.IsAllUpper("AB1")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
	assert.Equal(t, result3, false)
}

// TestIsAllLower 检查字符串是否全部由小写字母 a-z 组成。
// func IsAllLower(str string) bool
func TestIsAllLower(t *testing.T) {
	result1 := validator.IsAllLower("abc")
	result2 := validator.IsAllLower("abC")
	result3 := validator.IsAllLower("ab1")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
	assert.Equal(t, result3, false)
}

// TestIsPrintable 检查字符串是否全部由可打印组成。
// func IsPrintable(str string) bool
func TestIsPrintable(t *testing.T) {
	result1 := validator.IsPrintable("ABC")
	result2 := validator.IsPrintable("{id: 123}")
	result3 := validator.IsPrintable("")
	result4 := validator.IsPrintable("😄")
	result5 := validator.IsPrintable("\u0000")

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, true, result2, "result2")
	assert.Equal(t, true, result3, "result3")
	assert.Equal(t, true, result4, "result4")
	assert.Equal(t, false, result5, "result5")
}

// TestIsBin 检查给定的字符串是否为有效的二进制值。
// func IsBin(str string) bool
func TestIsBin(t *testing.T) {
	result1 := validator.IsBin("0101")
	result2 := validator.IsBin("0b1101")
	result3 := validator.IsBin("b1101")
	result4 := validator.IsBin("1201")

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, true, result2, "result2")
	assert.Equal(t, false, result3, "result3")
	assert.Equal(t, false, result4, "result4")
}

// TestIsHex 检查给定的字符串是否为有效的十六进制值。
// func IsHex(str string) bool
func TestIsHex(t *testing.T) {
	result1 := validator.IsHex("0xabcde")
	result2 := validator.IsHex("0XABCDE")
	result3 := validator.IsHex("cdfeg")
	result4 := validator.IsHex("0xcdfeg")

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, true, result2, "result2")
	assert.Equal(t, false, result3, "result3")
	assert.Equal(t, false, result4, "result4")
}

// TestIsASCII 检查字符串是否全部是ASCII字符。
// func IsIsASCIIAllLower(str string) bool
func TestIsASCII(t *testing.T) {
	result1 := validator.IsASCII("ABC")
	result2 := validator.IsASCII("123")
	result3 := validator.IsASCII("")
	result4 := validator.IsASCII("😄")
	result5 := validator.IsASCII("你好")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, true)
	assert.Equal(t, result3, true)
	assert.Equal(t, result4, false)
	assert.Equal(t, result5, false)
}

// TestIsBase64 检查字符串是否为base64字符串。
// func IsBase64(str string) bool
func TestIsBase64(t *testing.T) {
	result1 := validator.IsBase64("aGVsbG8=")
	result2 := validator.IsBase64("123456")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}

// TestIsBase64URL 检查给定的字符串是否为有效的URL - 安全的Base64编码字符串。
// func IsBase64URL(str string) bool
func TestIsBase64URL(t *testing.T) {
	result1 := validator.IsBase64URL("SAGsbG8sIHdvcmxkIQ")
	result2 := validator.IsBase64URL("SAGsbG8sIHdvcmxkIQ==")
	result3 := validator.IsBase64URL("SAGsbG8sIHdvcmxkIQ=")
	result4 := validator.IsBase64URL("SAGsbG8sIHdvcmxkIQ===")

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, false, result3)
	assert.Equal(t, false, result4)
}

// TestIsGBK 检查数据编码是否为 GBK（汉字内码扩展规范）。
// 此功能通过判断双字节是否在 GBK 的编码范围内来实现，而 UTF - 8 编码格式中的每个字节都在 GBK 的编码范围内。
// 因此，应首先调用 utf8.valid() 来检查是否不是 UTF - 8 编码，然后再调用 IsGBK() 来检查 GBK 编码
// func IsGBK(data []byte) bool
func TestIsGBK(t *testing.T) {
	str := "你好"
	gbkData, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(str))

	result := validator.IsGBK(gbkData)

	assert.Equal(t, result, true)
}

// TestIsJWT 检查给定的字符串是否是有效的JSON Web Token（JWT）。
// func IsJWT(str string) bool
func TestIsJWT(t *testing.T) {
	result1 := validator.IsJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibWVzc2FnZSI6IlB1dGluIGlzIGFic29sdXRlIHNoaXQiLCJpYXQiOjE1MTYyMzkwMjJ9.wkLWA5GtCpWdxNOrRse8yHZgORDgf8TpJp73WUQb910")
	result2 := validator.IsJWT("abc")

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, false, result2, "result2")
}
