package slice

import (
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/samber/lo"
	"github.com/samber/lo/parallel"

	"github.com/araujo88/lambda-go/pkg/core"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestMap 函数对片段中的每个元素应用给定的函数，返回一个包含转换后元素的新片段。
func TestMap(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		nums := []int{1, 2, 3}

		addOne := func(_ int, v int) int {
			return v + 1
		}

		result := slice.Map(nums, addOne)
		assert.Equal(t, []int{2, 3, 4}, result)
	}

	{
		numbers := []int{1, 2, 3, 4, 5}
		doubled := core.Map(numbers, func(x int) int { return x * 2 })
		assert.Equal(t, doubled, []int{2, 4, 6, 8, 10})
	}

	{
		result1 := lo.Map([]int{1, 2, 3, 4}, func(x int, _ int) string {
			return "Hello"
		})
		result2 := lo.Map([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
			return strconv.FormatInt(x, 10)
		})

		is.Equal(len(result1), 4)
		is.Equal(len(result2), 4)
		is.Equal(result1, []string{"Hello", "Hello", "Hello", "Hello"})
		is.Equal(result2, []string{"1", "2", "3", "4"})
	}
}

// TestMapConcurrent Applies the iteratee function to each item in the slice by concrrent.
// func MapConcurrent[T any, U any](slice []T, iteratee func(index int, item T) U, numThreads int) []U
func TestMapConcurrent(t *testing.T) {
	is := assert.New(t)

	{
		result1 := parallel.Map([]int{1, 2, 3, 4}, func(x int, _ int) string {
			return "Hello"
		})
		result2 := parallel.Map([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
			return strconv.FormatInt(x, 10)
		})

		is.Equal(len(result1), 4)
		is.Equal(len(result2), 4)
		is.Equal(result1, []string{"Hello", "Hello", "Hello", "Hello"})
		is.Equal(result2, []string{"1", "2", "3", "4"})
	}

	{
		nums := []int{1, 2, 3, 4, 5, 6}

		result := slice.MapConcurrent(nums, func(_, n int) int { return n * n }, 4)
		assert.Equal(t, []int{1, 4, 9, 16, 25, 36}, result)
	}

}

// TestForEach Iterates over elements of slice and invokes function for each element.
// func ForEach[T any](slice []T, iteratee func(index int, item T))
func TestForEach(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		// check of callback is called for every element and in proper order
		callParams1 := []string{}
		callParams2 := []int{}

		lo.ForEach([]string{"a", "b", "c"}, func(item string, i int) {
			callParams1 = append(callParams1, item)
			callParams2 = append(callParams2, i)
		})

		is.ElementsMatch([]string{"a", "b", "c"}, callParams1)
		is.ElementsMatch([]int{0, 1, 2}, callParams2)
		is.IsIncreasing(callParams2)
	}

	{
		nums := []int{1, 2, 3}

		var result []int
		addOne := func(_ int, v int) {
			result = append(result, v+1)
		}

		slice.ForEach(nums, addOne)
		assert.Equal(t, []int{2, 3, 4}, result)
	}
}

// TestForEachConcurrent Applies the iteratee function to each item in the slice concurrently.
// func ForEachConcurrent[T any](slice []T, iteratee func(index int, item T), numThreads int)
func TestForEachConcurrent(t *testing.T) {
	{
		nums := []int{1, 2, 3, 4, 5, 6, 7, 8}

		result := make([]int, len(nums))

		addOne := func(index int, value int) {
			result[index] = value + 1
		}

		slice.ForEachConcurrent(nums, addOne, 4)
		assert.Equal(t, []int{2, 3, 4, 5, 6, 7, 8, 9}, result)
	}

	{
		is := assert.New(t)

		var counter uint64
		collection := []int{1, 2, 3, 4}
		parallel.ForEach(collection, func(x int, i int) {
			atomic.AddUint64(&counter, 1)
		})

		is.Equal(uint64(4), atomic.LoadUint64(&counter))
	}
}

// TestForEachWithBreak Iterates over elements of slice and invokes function for each element, when iteratee return false, will break the for each loop.
// func ForEachWithBreak[T any](slice []T, iteratee func(index int, item T) bool)
func TestForEachWithBreak(t *testing.T) {
	{
		t.Parallel()
		is := assert.New(t)

		// check of callback is called for every element and in proper order

		var callParams1 []string
		var callParams2 []int

		lo.ForEachWhile([]string{"a", "b", "c"}, func(item string, i int) bool {
			if item == "c" {
				return false
			}
			callParams1 = append(callParams1, item)
			callParams2 = append(callParams2, i)
			return true
		})

		is.ElementsMatch([]string{"a", "b"}, callParams1)
		is.ElementsMatch([]int{0, 1}, callParams2)
		is.IsIncreasing(callParams2)
	}

	{
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
}
