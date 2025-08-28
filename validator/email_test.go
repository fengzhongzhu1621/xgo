package validator

import (
	"testing"

	"github.com/duke-git/lancet/v2/validator"
	"github.com/stretchr/testify/assert"
)

// 验证字符串是否是有效电子邮件地址
func TestIsEmail(t *testing.T) {
	result1 := validator.IsEmail("xx@xx.com")
	result2 := validator.IsEmail("x.x@@com")

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}
