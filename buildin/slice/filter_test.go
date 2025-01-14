package slice

import (
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestEvery Return all elements which match the function.
// func Filter[T any](slice []T, predicate func(index int, item T) bool) []T
func TestFilter(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	isEven := func(i, num int) bool {
		return num%2 == 0
	}

	result := slice.Filter(nums, isEven)

	assert.Equal(t, []int{2, 4}, result)
}

// TestFilterConcurrent Applies the provided filter function `predicate` to each element of the input slice concurrently.
// func FilterConcurrent[T any](slice []T, predicate func(index int, item T) bool, numThreads int) []T
func TestFilterConcurrent(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	isEven := func(i, num int) bool {
		return num%2 == 0
	}

	result := slice.FilterConcurrent(nums, isEven, 2)

	assert.Equal(t, []int{2, 4}, result)
}

// TestFilterMap Returns a slice which apply both filtering and mapping to the given slice. iteratee callback function should returntwo values: 1, mapping result. 2, whether the result element should be included or not.
// func FilterMap[T any, U any](slice []T, iteratee func(index int, item T) (U, bool)) []U
func TestFilterMap(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	getEvenNumStr := func(i, num int) (string, bool) {
		if num%2 == 0 {
			return strconv.FormatInt(int64(num), 10), true
		}
		return "", false
	}

	result := slice.FilterMap(nums, getEvenNumStr)
	assert.Equal(t, []string{"2", "4"}, result)
}
