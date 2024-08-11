package validator

import (
	"testing"

	"github.com/duke-git/lancet/v2/validator"
	"github.com/stretchr/testify/assert"
)

// 验证字符串是否包含中文字符
func TestContainChinese(t *testing.T) {
	result1 := validator.ContainChinese("你好")
	result2 := validator.ContainChinese("你好 hello")
	result3 := validator.ContainChinese("hello")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, true)
	assert.Equal(t, result3, false)
}

// 验证字符串是否是中国手机号码
func TestIsChineseMobile(t *testing.T) {
	result1 := validator.IsChineseMobile("17530367777")
	result2 := validator.IsChineseMobile("121212121")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}

// 验证字符串是否是中国身份证号码
func TestIsChineseIdNum(t *testing.T) {
	result1 := validator.IsChineseIdNum("210911192105130715")
	result2 := validator.IsChineseIdNum("123456")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}

// 验证字符串是否是中国电话座机号码
func TestIsChinesePhone(t *testing.T) {
	result1 := validator.IsChinesePhone("010-32116675")
	result2 := validator.IsChinesePhone("123-87562")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}

// 验证字符串是否是信用卡号码
func TestIsCreditCard(t *testing.T) {
	result1 := validator.IsCreditCard("4111111111111111")
	result2 := validator.IsCreditCard("123456")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}
