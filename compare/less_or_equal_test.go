package compare

import (
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/compare"
	"github.com/stretchr/testify/assert"
)

// TestLessOrEqual godoc
// func LessOrEqual(left, right any) bool
func TestLessOrEqual(t *testing.T) {
	result1 := compare.LessOrEqual(1, 1)
	result2 := compare.LessOrEqual(1.1, 2.2)
	result3 := compare.LessOrEqual("a", "b")

	time1 := time.Now()
	time2 := time1.Add(time.Second)
	result4 := compare.LessOrEqual(time1, time2)

	result5 := compare.LessOrEqual(2, 1)
	result6 := compare.LessOrEqual(1, int64(2))

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, true, result3)
	assert.Equal(t, true, result4)
	assert.Equal(t, false, result5)
	assert.Equal(t, false, result6)
}
