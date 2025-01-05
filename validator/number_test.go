package validator

import (
	"testing"

	"github.com/duke-git/lancet/v2/validator"
	"github.com/stretchr/testify/assert"
)

// TestIsInt 检查值是否为整数（int，unit）或不是。
// func IsInt(v any) bool
func TestIsInt(t *testing.T) {
	result1 := validator.IsInt("")
	result2 := validator.IsInt("3")
	result3 := validator.IsInt(0.1)
	result4 := validator.IsInt(0)

	assert.Equal(t, false, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, false, result3)
	assert.Equal(t, true, result4)
}

// TestIsFloat 检查值是否为浮点数(float32, float34) 或不是。
// func IsFloat(v any) bool
func TestIsFloat(t *testing.T) {
	result1 := validator.IsFloat("")
	result2 := validator.IsFloat("3")
	result3 := validator.IsFloat(0)
	result4 := validator.IsFloat(0.1)

	assert.Equal(t, false, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, false, result3)
	assert.Equal(t, true, result4)
}

// TestIsNumber 检查值是否为数字(integer, float) 或不是。
// func IsNumber(v any) bool
func TestIsNumber(t *testing.T) {
	result1 := validator.IsNumber("")
	result2 := validator.IsNumber("3")
	result3 := validator.IsNumber(0.1)
	result4 := validator.IsNumber(0)
	result5 := validator.IsNumber(int64(1))

	assert.Equal(t, false, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, true, result3)
	assert.Equal(t, true, result4)
	assert.Equal(t, true, result5)
}

// TestIsIntStr 检查字符串是否可以转换为整数。
// func IsIntStr(v any) bool
func TestIsIntStr(t *testing.T) {
	result1 := validator.IsIntStr("+3")
	result2 := validator.IsIntStr("-3")
	result3 := validator.IsIntStr("3.")
	result4 := validator.IsIntStr("abc")

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, false, result3)
	assert.Equal(t, false, result4)
}

// TestIsFloatStr 检查字符串是否可以转换为浮点数。
// func IsFloatStr(v any) bool
func TestIsFloatStr(t *testing.T) {
	result1 := validator.IsFloatStr("3.")
	result2 := validator.IsFloatStr("+3.")
	result3 := validator.IsFloatStr("12")
	result4 := validator.IsFloatStr("abc")

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, true, result3)
	assert.Equal(t, false, result4)
}

// TestIsNumberStr 检查字符串是否可以转换为数值。
// func IsNumberStr(v any) bool
func TestIsNumberStr(t *testing.T) {
	result1 := validator.IsNumberStr("3.")
	result2 := validator.IsNumberStr("+3.")
	result3 := validator.IsNumberStr("+3e2")
	result4 := validator.IsNumberStr("abc")

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, true, result3)
	assert.Equal(t, false, result4)
}
