package compare

import (
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/compare"
	"github.com/stretchr/testify/assert"
)

// TestLessThan检查两个值是否相等。（仅检查值）
// func LessThan(left, right any) bool
func TestLessThan(t *testing.T) {
	result1 := compare.LessThan(1, 2)
	result2 := compare.LessThan(1.1, 2.2)
	result3 := compare.LessThan("a", "b")

	time1 := time.Now()
	time2 := time1.Add(time.Second)
	result4 := compare.LessThan(time1, time2)

	result5 := compare.LessThan(2, 1)
	result6 := compare.LessThan(1, int64(2))

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, true, result3)
	assert.Equal(t, true, result4)
	assert.Equal(t, false, result5)
	assert.Equal(t, false, result6)
}
