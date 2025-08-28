package validator

import (
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestIsAscending Checks if a slice is ascending order.
// func IsAscending[T constraints.Ordered](slice []T) bool
func TestIsAscending(t *testing.T) {
	result1 := slice.IsAscending([]int{1, 2, 3, 4, 5})
	result2 := slice.IsAscending([]int{5, 4, 3, 2, 1})
	result3 := slice.IsAscending([]int{2, 1, 3, 4, 5})

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, false, result3)
}

// TestIsDescending Checks if a slice is descending order.
// func IsDescending[T constraints.Ordered](slice []T) bool
func TestIsDescending(t *testing.T) {
	result1 := slice.IsDescending([]int{5, 4, 3, 2, 1})
	result2 := slice.IsDescending([]int{1, 2, 3, 4, 5})
	result3 := slice.IsDescending([]int{2, 1, 3, 4, 5})

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, false, result3)
}

// TestIsSorted Checks if a slice is sorted (ascending or descending).
// func IsSorted[T constraints.Ordered](slice []T) bool
func TestIsSorted(t *testing.T) {
	{
		result1 := slice.IsSorted([]int{5, 4, 3, 2, 1})
		result2 := slice.IsSorted([]int{1, 2, 3, 4, 5})
		result3 := slice.IsSorted([]int{2, 1, 3, 4, 5})

		assert.Equal(t, true, result1)
		assert.Equal(t, true, result2)
		assert.Equal(t, false, result3)
	}

	{
		t.Parallel()
		is := assert.New(t)

		is.True(lo.IsSortedByKey([]string{"a", "bb", "ccc"}, func(s string) int {
			return len(s)
		}))

		is.False(lo.IsSortedByKey([]string{"aa", "b", "ccc"}, func(s string) int {
			return len(s)
		}))

		is.True(lo.IsSortedByKey([]string{"1", "2", "3", "11"}, func(s string) int {
			ret, _ := strconv.Atoi(s)
			return ret
		}))
	}
}

// TestIsSortedByKey Checks if a slice is sorted by iteratee function.
// func IsSortedByKey[T any, K constraints.Ordered](slice []T, iteratee func(item T) K) bool
func TestIsSortedByKey(t *testing.T) {
	result1 := slice.IsSortedByKey([]string{"a", "ab", "abc"}, func(s string) int {
		return len(s)
	})
	result2 := slice.IsSortedByKey([]string{"abc", "ab", "a"}, func(s string) int {
		return len(s)
	})
	result3 := slice.IsSortedByKey([]string{"abc", "a", "ab"}, func(s string) int {
		return len(s)
	})

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, false, result3)
}
