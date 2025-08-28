package validator

import (
	"testing"

	"github.com/duke-git/lancet/v2/validator"
	"github.com/stretchr/testify/assert"
)

// TestIsChineseIdNum 验证字符串是否是中国身份证号码
func TestIsChineseIdNum(t *testing.T) {
	result1 := validator.IsChineseIdNum("210911192105130715")
	result2 := validator.IsChineseIdNum("123456")

	assert.Equal(t, false, result1, "result1")
	assert.Equal(t, false, result2, "result2")
}

// TestIsCreditCard 验证字符串是否是信用卡号码
func TestIsCreditCard(t *testing.T) {
	result1 := validator.IsCreditCard("4111111111111111")
	result2 := validator.IsCreditCard("123456")

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}

// TestIsVisa 检查所给的字符串是否为有效的Visa卡号。
// func IsVisa(v string) bool
func TestIsVisa(t *testing.T) {
	result1 := validator.IsVisa("4111111111111111")
	result2 := validator.IsVisa("123")

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}

// TestIsMasterCard 检查所给的字符串是否为有效的Master卡号。
// func IsMasterCard(v string) bool
func TestIsMasterCard(t *testing.T) {
	result1 := validator.IsMasterCard("5425233430109903")
	result2 := validator.IsMasterCard("4111111111111111")

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}

// TestIsAmericanExpress 检查所给的字符串是否为有效的美国运通卡号。
// func IsAmericanExpress(v string) bool
func TestIsAmericanExpress(t *testing.T) {
	result1 := validator.IsAmericanExpress("342883359122187")
	result2 := validator.IsAmericanExpress("3782822463100007")

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}

// TestIsUnionPay 检查所给的字符串是否为有效的银联卡号。
// func func IsUnionPay(v string) bool
func TestIsUnionPay(t *testing.T) {
	result1 := validator.IsUnionPay("6221263430109903")
	result2 := validator.IsUnionPay("3782822463100007")

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}

// TestIsChinaUnionPay 检查所给的字符串是否为有效的中国银联卡号。
// func func IsChinaUnionPay(v string) bool
func TestIsChinaUnionPay(t *testing.T) {
	result1 := validator.IsChinaUnionPay("6250941006528599")
	result2 := validator.IsChinaUnionPay("3782822463100007")

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}
