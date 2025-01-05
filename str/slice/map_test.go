package slice

import (
	"strconv"
	"testing"

	"github.com/araujo88/lambda-go/pkg/core"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestMap 函数对片段中的每个元素应用给定的函数，返回一个包含转换后元素的新片段。
func TestMap(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	doubled := core.Map(numbers, func(x int) int { return x * 2 })
	assert.Equal(t, doubled, []int{2, 4, 6, 8, 10})
}

// TestLancetMap Creates an slice of values by running each element in slice thru function.
// func Map[T any, U any](slice []T, iteratee func(index int, item T) U) []U
func TestLancetMap(t *testing.T) {
	nums := []int{1, 2, 3}

	addOne := func(_ int, v int) int {
		return v + 1
	}

	result := slice.Map(nums, addOne)
	assert.Equal(t, []int{2, 3, 4}, result)
}

// TestMapConcurrent Applies the iteratee function to each item in the slice by concrrent.
// func MapConcurrent[T any, U any](slice []T, iteratee func(index int, item T) U, numThreads int) []U
func TestMapConcurrent(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}

	result := slice.MapConcurrent(nums, func(_, n int) int { return n * n }, 4)
	assert.Equal(t, []int{1, 4, 9, 16, 25, 36}, result)
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

// TestForEach Iterates over elements of slice and invokes function for each element.
// func ForEach[T any](slice []T, iteratee func(index int, item T))
func TestForEach(t *testing.T) {
	nums := []int{1, 2, 3}

	var result []int
	addOne := func(_ int, v int) {
		result = append(result, v+1)
	}

	slice.ForEach(nums, addOne)
	assert.Equal(t, []int{2, 3, 4}, result)
}

// TestForEachConcurrent Applies the iteratee function to each item in the slice concurrently.
// func ForEachConcurrent[T any](slice []T, iteratee func(index int, item T), numThreads int)
func TestForEachConcurrent(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}

	result := make([]int, len(nums))

	addOne := func(index int, value int) {
		result[index] = value + 1
	}

	slice.ForEachConcurrent(nums, addOne, 4)
	assert.Equal(t, []int{2, 3, 4, 5, 6, 7, 8, 9}, result)
}

// TestForEachWithBreak Iterates over elements of slice and invokes function for each element, when iteratee return false, will break the for each loop.
// func ForEachWithBreak[T any](slice []T, iteratee func(index int, item T) bool)
func TestForEachWithBreak(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	var sum int

	slice.ForEachWithBreak(numbers, func(_, n int) bool {
		if n > 3 {
			// 停止循环
			return false
		}
		sum += n
		return true
	})

	// 1 + 2 + 3
	assert.Equal(t, 6, sum)
}
