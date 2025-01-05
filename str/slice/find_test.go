package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestFindBy Iterates over elements of slice, returning the first one that passes a truth test on predicate function.If return T is nil or zero value then no items matched the predicate func. In contrast to Find or FindLast, its return value no longer requires dereferencing.
// 遍历切片中的元素，返回第一个通过谓词函数真值测试的元素。如果返回值 T 是 nil 或零值，那么没有元素匹配该谓词函数。与 Find 或 FindLast 不同，其返回值不再需要解引用。
// func FindBy[T any](slice []T, predicate func(index int, item T) bool) (v T, ok bool)
func TestFindBy(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	isEven := func(i, num int) bool {
		return num%2 == 0
	}

	result, ok := slice.FindBy(nums, isEven)

	assert.Equal(t, 2, result)
	assert.Equal(t, true, ok)
}

// TestFindLastBy FindLastBy iterates over elements of slice, returning the last one that passes a truth test on predicate function. If return T is nil or zero value then no items matched the predicate func. In contrast to Find or FindLast, its return value no longer requires dereferencing.
// 遍历切片中的元素，返回第一个通过谓词函数真值测试的元素。如果返回值 T 是 nil 或零值，那么没有元素匹配该谓词函数。与 Find 或 FindLast 不同，其返回值不再需要解引用。
// func FindLastBy[T any](slice []T, predicate func(index int, item T) bool) (v T, ok bool)
func TestFindLastBy(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	isEven := func(i, num int) bool {
		return num%2 == 0
	}

	result, ok := slice.FindLastBy(nums, isEven)

	assert.Equal(t, 4, result)
	assert.Equal(t, true, ok)
}

// TestIndexOf Returns the index at which the first occurrence of a item is found in a slice or return -1 if the item cannot be found.
// func IndexOf[T comparable](slice []T, item T) int
func TestIndexOf(t *testing.T) {
	strs := []string{"a", "a", "b", "c"}

	result1 := slice.IndexOf(strs, "a")
	result2 := slice.IndexOf(strs, "d")

	assert.Equal(t, 0, result1)
	assert.Equal(t, -1, result2)
}

// TestLastIndexOf Returns the index at which the last occurrence of a item is found in a slice or return -1 if the item cannot be found.
// func LastIndexOf[T comparable](slice []T, item T) int
func TestLastIndexOf(t *testing.T) {
	strs := []string{"a", "a", "b", "c"}

	result1 := slice.LastIndexOf(strs, "a")
	result2 := slice.LastIndexOf(strs, "d")

	assert.Equal(t, 1, result1)
	assert.Equal(t, -1, result2)
}
