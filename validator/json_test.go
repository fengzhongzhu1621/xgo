package validator

import (
	"testing"

	"github.com/duke-git/lancet/v2/validator"
	"github.com/stretchr/testify/assert"
)

// TestIsJSON 检查字符串是否是有效的JSON。
// func IsJSON(v any) bool
func TestIsJSON(t *testing.T) {
	result1 := validator.IsJSON("{}")
	result2 := validator.IsJSON("{\"name\": \"test\"}")
	result3 := validator.IsJSON("")
	result4 := validator.IsJSON("abc")

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, false, result3)
	assert.Equal(t, false, result4)
}
