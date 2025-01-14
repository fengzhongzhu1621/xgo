package slice

import (
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestFlatten Flatten slice with one level.
// func Flatten(slice any) any
func TestFlatten(t *testing.T) {
	arrs := [][][]string{{{"a", "b"}}, {{"c", "d"}}}

	result1 := slice.Flatten(arrs)
	assert.Equal(t, [][]string{{"a", "b"}, {"c", "d"}}, result1)

	result2 := slice.Flatten(result1)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result2)
}

// TestFlattenDeep flattens slice recursive.
// func FlattenDeep(slice any) any
func TestFlattenDeep(t *testing.T) {
	arrs := [][][]string{{{"a", "b"}}, {{"c", "d"}}}

	result := slice.FlattenDeep(arrs)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result)
}

// TestFlatMap Manipulates a slice and transforms and flattens it to a slice of another type.
// func FlatMap[T any, U any](slice []T, iteratee func(index int, item T) []U) []U
func TestFlatMap(t *testing.T) {
	nums := []int{1, 2, 3, 4}

	result := slice.FlatMap(nums, func(i int, num int) []string {
		s := "hi-" + strconv.FormatInt(int64(num), 10)
		return []string{s}
	})

	assert.Equal(t, []string{"hi-1", "hi-2", "hi-3", "hi-4"}, result)
}
