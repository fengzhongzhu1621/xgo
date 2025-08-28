package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestCount Returns the number of occurrences of the given item in the slice.
// func Count[T comparable](slice []T, item T) int
func TestCount(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		count1 := lo.Count([]int{1, 2, 1}, 1)
		count2 := lo.Count([]int{1, 2, 1}, 3)
		count3 := lo.Count([]int{}, 1)

		is.Equal(count1, 2)
		is.Equal(count2, 0)
		is.Equal(count3, 0)
	}

	{
		nums := []int{1, 2, 3, 3, 4}

		result1 := slice.Count(nums, 1)
		result2 := slice.Count(nums, 3)

		assert.Equal(t, 1, result1)
		assert.Equal(t, 2, result2)
	}
}

// TestCountBy Iterates over elements of slice with predicate function, returns the number of all matched elements.
// func CountBy[T any](slice []T, predicate func(index int, item T) bool) int
func TestCountBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		count1 := lo.CountBy([]int{1, 2, 1}, func(i int) bool {
			return i < 2
		})

		count2 := lo.CountBy([]int{1, 2, 1}, func(i int) bool {
			return i > 2
		})

		count3 := lo.CountBy([]int{}, func(i int) bool {
			return i <= 2
		})

		is.Equal(count1, 2)
		is.Equal(count2, 0)
		is.Equal(count3, 0)
	}

	{
		nums := []int{1, 2, 3, 4, 5}

		isEven := func(i, num int) bool {
			return num%2 == 0
		}

		result := slice.CountBy(nums, isEven)
		assert.Equal(t, 1, result)
	}
}

// TestFrequency Counts the frequency of each element in the slice.
// func Frequency[T comparable](slice []T) map[T]int
func TestFrequency(t *testing.T) {
	{
		strs := []string{"a", "b", "b", "c", "c", "c"}
		result := slice.Frequency(strs)

		assert.Equal(t, map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}, result)
	}

	{
		t.Parallel()
		is := assert.New(t)

		is.Equal(map[int]int{}, lo.CountValues([]int{}))
		is.Equal(map[int]int{1: 1, 2: 1}, lo.CountValues([]int{1, 2}))
		is.Equal(map[int]int{1: 1, 2: 2}, lo.CountValues([]int{1, 2, 2}))
		is.Equal(
			map[string]int{"": 1, "foo": 1, "bar": 1},
			lo.CountValues([]string{"foo", "bar", ""}),
		)
		is.Equal(map[string]int{"foo": 1, "bar": 2}, lo.CountValues([]string{"foo", "bar", "bar"}))
	}
}
