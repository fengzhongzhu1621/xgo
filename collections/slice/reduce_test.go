package slice

import (
	"fmt"
	"testing"

	"github.com/araujo88/lambda-go/pkg/core"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// Foldl 和 Foldr 是两个功能强大的函数，通过对每个元素应用一个函数和一个累加器，可以将切片还原为一个单一值。
// 这两个函数的区别在于它们遍历切片的方向。
func TestFold(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	sum := core.Foldl(func(acc, x int) int { return acc + x }, 0, numbers)
	assert.Equal(t, 15, sum)
}

// TestReduceBy Produces a value from slice by accumulating the result of each element as passed through the reducer function.
// func ReduceBy[T any, U any](slice []T, initial U, reducer func(index int, item T, agg U) U) U
func TestReduceBy(t *testing.T) {
	result1 := slice.ReduceBy([]int{1, 2, 3, 4}, 0, func(_ int, item int, agg int) int {
		return agg + item
	})

	result2 := slice.ReduceBy([]int{1, 2, 3, 4}, "", func(_ int, item int, agg string) string {
		return agg + fmt.Sprintf("%v", item)
	})

	assert.Equal(t, 10, result1)
	assert.Equal(t, 1234, result2)
}

// TestReduceConcurrent Reduces the slice to a single value by applying the reducer function to each item in the slice concurrently.
// func ReduceConcurrent[T any](slice []T, initial T, reducer func(index int, item T, agg T) T, numThreads int) T
func TestReduceConcurrent(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := slice.ReduceConcurrent(nums, 0, func(_ int, item, agg int) int {
		return agg + item
	}, 1)

	assert.Equal(t, 55, result)
}

// TestReduceRight ReduceRight is like ReduceBy, but it iterates over elements of slice from right to left.
// func ReduceRight[T any, U any](slice []T, initial U, reducer func(index int, item T, agg U) U) U
func TestReduceRight(t *testing.T) {
	result := slice.ReduceRight([]int{1, 2, 3, 4}, "", func(_ int, item int, agg string) string {
		return agg + fmt.Sprintf("%v", item)
	})

	assert.Equal(t, "4321", result)
}
