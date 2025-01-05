package validator

import (
	"testing"

	"github.com/duke-git/lancet/v2/validator"
	"github.com/stretchr/testify/assert"
)

// TestIsZeroValue 检查字符串是否是零值
// func IsZeroValue(value any) bool
func TestIsZeroValue(t *testing.T) {
	result1 := validator.IsZeroValue("")
	result2 := validator.IsZeroValue(0)
	result3 := validator.IsZeroValue("abc")
	result4 := validator.IsZeroValue(1)

	assert.Equal(t, true, result1, "result1")
	assert.Equal(t, true, result2, "result2")
	assert.Equal(t, false, result3, "result3")
	assert.Equal(t, false, result4, "result4")
}
