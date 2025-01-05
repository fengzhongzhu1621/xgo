package validator

import (
	"testing"

	"github.com/duke-git/lancet/v2/validator"
	"github.com/stretchr/testify/assert"
)

// TestIsStrongPassword 检查字符串是否是强密码（小写和大写字母+数字+特殊字符（! @ # $ % ^ & *（）？> <
// func IsStrongPassword(password string, length int) bool
func TestIsStrongPassword(t *testing.T) {
	result1 := validator.IsStrongPassword("abcABC", 6)
	result2 := validator.IsStrongPassword("abcABC123@#$", 10)

	assert.Equal(t, false, result1, "result1")
	assert.Equal(t, true, result2, "result2")
}

// TestIsWeakPassword 检查字符串是否为弱密码（仅字母或仅数字或字母+数字）。
// func IsWeakPassword(password string, length int) bool
func TestIsWeakPassword(t *testing.T) {
	result1 := validator.IsWeakPassword("abcABC")
	result2 := validator.IsWeakPassword("abc123@#$")

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, false, result2, "result2")
}
