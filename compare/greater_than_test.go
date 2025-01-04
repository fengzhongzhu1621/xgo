package compare

import (
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/compare"
	"github.com/stretchr/testify/assert"
)

// TestGreaterThan godoc
// func GreaterThan(left, right any) bool
func TestGreaterThan(t *testing.T) {
	result1 := compare.GreaterThan(2, 1)
	result2 := compare.GreaterThan(2.2, 1.1)
	result3 := compare.GreaterThan("b", "a")

	time1 := time.Now()
	time2 := time1.Add(time.Second)
	result4 := compare.GreaterThan(time2, time1)

	result5 := compare.GreaterThan(1, 2)
	result6 := compare.GreaterThan(int64(2), 1)
	result7 := compare.GreaterThan("b", "c")

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, true, result3)
	assert.Equal(t, true, result4)
	assert.Equal(t, false, result5)
	assert.Equal(t, false, result6)
	assert.Equal(t, false, result7)
}
