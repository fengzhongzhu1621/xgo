package compare

import (
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/compare"
	"github.com/stretchr/testify/assert"
)

// TestGreaterOrEqual godoc
// func TestGreaterOrEqual(left, right any) bool
func TestGreaterOrEqual(t *testing.T) {
	result1 := compare.GreaterOrEqual(1, 1)
	result2 := compare.GreaterOrEqual(2.2, 1.1)
	result3 := compare.GreaterOrEqual("b", "b")

	time1 := time.Now()
	time2 := time1.Add(time.Second)
	result4 := compare.GreaterOrEqual(time2, time1)

	result5 := compare.GreaterOrEqual(1, 2)
	result6 := compare.GreaterOrEqual(int64(2), 1)
	result7 := compare.GreaterOrEqual("b", "c")

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, true, result3)
	assert.Equal(t, true, result4)
	assert.Equal(t, false, result5)
	assert.Equal(t, false, result6)
	assert.Equal(t, false, result7)
}
