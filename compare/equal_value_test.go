package compare

import (
	"testing"

	"github.com/duke-git/lancet/v2/compare"
	"github.com/stretchr/testify/assert"
)

// TestEqualValue 检查两个值是否相等。（仅检查值）
// Checks if two values are equal or not. (check value only)
func TestEqualValue(t *testing.T) {
	result1 := compare.EqualValue(1, 1)
	result2 := compare.EqualValue(int(1), int64(1))
	result3 := compare.EqualValue(1, "1")
	result4 := compare.EqualValue(1, "2")

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, true, result3)
	assert.Equal(t, false, result4)
}
