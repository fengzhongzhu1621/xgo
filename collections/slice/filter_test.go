package slice

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/samber/lo"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestEvery Return all elements which match the function.
// func Filter[T any](slice []T, predicate func(index int, item T) bool) []T
func TestFilter(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		r1 := lo.Filter([]int{1, 2, 3, 4}, func(x int, _ int) bool {
			return x%2 == 0
		})
		is.Equal(r1, []int{2, 4})

		r2 := lo.Filter([]string{"", "foo", "", "bar", ""}, func(x string, _ int) bool {
			return len(x) > 0
		})
		is.Equal(r2, []string{"foo", "bar"})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Filter(allStrings, func(x string, _ int) bool {
			return len(x) > 0
		})
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
		nums := []int{1, 2, 3, 4, 5}
		isEven := func(i, num int) bool {
			return num%2 == 0
		}
		result := slice.Filter(nums, isEven)
		assert.Equal(t, []int{2, 4}, result)
	}
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
	t.Parallel()
	is := assert.New(t)

	{
		r1 := lo.FilterMap([]int64{1, 2, 3, 4}, func(x int64, _ int) (string, bool) {
			if x%2 == 0 {
				return strconv.FormatInt(x, 10), true
			}
			return "", false
		})
		r2 := lo.FilterMap([]string{"cpu", "gpu", "mouse", "keyboard"}, func(x string, _ int) (string, bool) {
			if strings.HasSuffix(x, "pu") {
				return "xpu", true
			}
			return "", false
		})

		is.Equal(len(r1), 2)
		is.Equal(len(r2), 2)
		is.Equal(r1, []string{"2", "4"})
		is.Equal(r2, []string{"xpu", "xpu"})
	}

	{
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
}

// TestSubset returns a copy of a slice from `offset` up to `length` elements. Like `slice[start:start+length]`, but does not panic on overflow.
func TestSubset(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	in := []int{0, 1, 2, 3, 4}

	out1 := lo.Subset(in, 0, 0)
	out2 := lo.Subset(in, 10, 2)
	out3 := lo.Subset(in, -10, 2)
	out4 := lo.Subset(in, 0, 10)
	out5 := lo.Subset(in, 0, 2)
	out6 := lo.Subset(in, 2, 2)
	out7 := lo.Subset(in, 2, 5)
	out8 := lo.Subset(in, 2, 3)
	out9 := lo.Subset(in, 2, 4)
	out10 := lo.Subset(in, -2, 4)
	out11 := lo.Subset(in, -4, 1)
	out12 := lo.Subset(in, -4, math.MaxUint)

	is.Equal([]int{}, out1)
	is.Equal([]int{}, out2)
	is.Equal([]int{0, 1}, out3)
	is.Equal([]int{0, 1, 2, 3, 4}, out4)
	is.Equal([]int{0, 1}, out5)
	is.Equal([]int{2, 3}, out6)
	is.Equal([]int{2, 3, 4}, out7)
	is.Equal([]int{2, 3, 4}, out8)
	is.Equal([]int{2, 3, 4}, out9)
	is.Equal([]int{3, 4}, out10)
	is.Equal([]int{1}, out11)
	is.Equal([]int{1, 2, 3, 4}, out12)

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.Subset(allStrings, 0, 2)
	is.IsType(nonempty, allStrings, "type preserved")
}

func TestSlice(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	in := []int{0, 1, 2, 3, 4}

	out1 := lo.Slice(in, 0, 0)
	out2 := lo.Slice(in, 0, 1)
	out3 := lo.Slice(in, 0, 5)
	out4 := lo.Slice(in, 0, 6)
	out5 := lo.Slice(in, 1, 1)
	out6 := lo.Slice(in, 1, 5)
	out7 := lo.Slice(in, 1, 6)
	out8 := lo.Slice(in, 4, 5)
	out9 := lo.Slice(in, 5, 5)
	out10 := lo.Slice(in, 6, 5)
	out11 := lo.Slice(in, 6, 6)
	out12 := lo.Slice(in, 1, 0)
	out13 := lo.Slice(in, 5, 0)
	out14 := lo.Slice(in, 6, 4)
	out15 := lo.Slice(in, 6, 7)
	out16 := lo.Slice(in, -10, 1)
	out17 := lo.Slice(in, -1, 3)
	out18 := lo.Slice(in, -10, 7)

	is.Equal([]int{}, out1)
	is.Equal([]int{0}, out2)
	is.Equal([]int{0, 1, 2, 3, 4}, out3)
	is.Equal([]int{0, 1, 2, 3, 4}, out4)
	is.Equal([]int{}, out5)
	is.Equal([]int{1, 2, 3, 4}, out6)
	is.Equal([]int{1, 2, 3, 4}, out7)
	is.Equal([]int{4}, out8)
	is.Equal([]int{}, out9)
	is.Equal([]int{}, out10)
	is.Equal([]int{}, out11)
	is.Equal([]int{}, out12)
	is.Equal([]int{}, out13)
	is.Equal([]int{}, out14)
	is.Equal([]int{}, out15)
	is.Equal([]int{0}, out16)
	is.Equal([]int{0, 1, 2}, out17)
	is.Equal([]int{0, 1, 2, 3, 4}, out18)

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.Slice(allStrings, 0, 2)
	is.IsType(nonempty, allStrings, "type preserved")
}

// TestWithout Creates a slice excluding all given items.
// func Without[T comparable](slice []T, items ...T) []T
func TestWithout(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result := slice.Without([]int{1, 2, 3, 4}, 1, 2)
		assert.Equal(t, []int{3, 4}, result)

	}

	{
		result1 := lo.Without([]int{0, 2, 10}, 0, 1, 2, 3, 4, 5)
		result2 := lo.Without([]int{0, 7}, 0, 1, 2, 3, 4, 5)
		result3 := lo.Without([]int{}, 0, 1, 2, 3, 4, 5)
		result4 := lo.Without([]int{0, 1, 2}, 0, 1, 2)
		result5 := lo.Without([]int{})
		is.Equal(result1, []int{10})
		is.Equal(result2, []int{7})
		is.Equal(result3, []int{})
		is.Equal(result4, []int{})
		is.Equal(result5, []int{})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Without(allStrings, "")
		is.IsType(nonempty, allStrings, "type preserved")
	}
}