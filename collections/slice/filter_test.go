package slice

import (
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
