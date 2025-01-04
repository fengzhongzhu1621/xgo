package compare

import (
	"testing"

	"github.com/duke-git/lancet/v2/compare"
	"github.com/stretchr/testify/assert"
)

// TestInDelta 检查两个值是否在增量范围内相等。
// func InDelta[T constraints.Integer | constraints.Float](left, right T, delta float64) bool
func TestInDelta(t *testing.T) {
	result1 := compare.InDelta(1, 1, 0)
	result2 := compare.InDelta(1, 2, 0)

	result3 := compare.InDelta(2.0/3.0, 0.66667, 0.001)
	result4 := compare.InDelta(2.0/3.0, 0.0, 0.001)

	result5 := compare.InDelta(float64(74.96)-float64(20.48), 54.48, 0)
	result6 := compare.InDelta(float64(74.96)-float64(20.48), 54.48, 1e-14)

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, true, result3)
	assert.Equal(t, false, result4)
	assert.Equal(t, false, result5)
	assert.Equal(t, true, result6)

}
