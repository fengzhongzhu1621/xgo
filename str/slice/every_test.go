package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestEvery Return true if all of the values in the slice pass the predicate function.
// func Every[T any](slice []T, predicate func(index int, item T) bool) bool
func TestEvery(t *testing.T) {
	nums := []int{1, 2, 3, 5}

	isEven := func(i, num int) bool {
		return num%2 == 0
	}

	result := slice.Every(nums, isEven)

	assert.Equal(t, false, result)
}
