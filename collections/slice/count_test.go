package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestCount Returns the number of occurrences of the given item in the slice.
// func Count[T comparable](slice []T, item T) int
func TestCount(t *testing.T) {
	nums := []int{1, 2, 3, 3, 4}

	result1 := slice.Count(nums, 1)
	result2 := slice.Count(nums, 3)

	assert.Equal(t, 1, result1)
	assert.Equal(t, 2, result2)
}

// TestCountBy Iterates over elements of slice with predicate function, returns the number of all matched elements.
// func CountBy[T any](slice []T, predicate func(index int, item T) bool) int
func TestCountBy(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	isEven := func(i, num int) bool {
		return num%2 == 0
	}

	result := slice.CountBy(nums, isEven)
	assert.Equal(t, 1, result)
}
