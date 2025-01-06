package slice

import (
	"github.com/duke-git/lancet/v2/slice"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSome Return true if any of the values in the list pass the predicate function.
// func Some[T any](slice []T, predicate func(index int, item T) bool) bool
func TestSome(t *testing.T) {
	nums := []int{1, 2, 3, 5}

	isEven := func(i, num int) bool {
		return num%2 == 0
	}

	result := slice.Some(nums, isEven)

	assert.Equal(t, true, result)
}
