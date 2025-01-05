package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestDeleteAt delete the element of slice at index.
// func DeleteAt[T any](slice []T, index int)
func TestDeleteAt(t *testing.T) {
	chars := []string{"a", "b", "c", "d", "e"}

	result1 := slice.DeleteAt(chars, 0)
	result2 := slice.DeleteAt(chars, 4)
	result3 := slice.DeleteAt(chars, 5)
	result4 := slice.DeleteAt(chars, 6)
	result5 := slice.DeleteAt(chars, 100)

	assert.Equal(t, []string{"b", "c", "d", "e"}, result1)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result2)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result3)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result4)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result5)
}

// TestDeleteRange Delete the element of slice from start index to end index（exclude)
// func DeleteRange[T any](slice []T, start, end int) []T
func TestDeleteRange(t *testing.T) {
	chars := []string{"a", "b", "c", "d", "e"}

	result1 := slice.DeleteRange(chars, 0, 0)
	result2 := slice.DeleteRange(chars, 0, 1)
	result3 := slice.DeleteRange(chars, 0, 3)
	result4 := slice.DeleteRange(chars, 0, 4)
	result5 := slice.DeleteRange(chars, 0, 5)

	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, result1)
	assert.Equal(t, []string{"b", "c", "d", "e"}, result2)
	assert.Equal(t, []string{"d", "e"}, result3)
	assert.Equal(t, []string{"e"}, result4)
	assert.Equal(t, []string{}, result5)
}

// TestDrop Drop n elements from the start of a slice.
// func Drop[T any](slice []T, n int) []T
func TestDrop(t *testing.T) {
	result1 := slice.Drop([]string{"a", "b", "c"}, 0)
	result2 := slice.Drop([]string{"a", "b", "c"}, 1)
	result3 := slice.Drop([]string{"a", "b", "c"}, -1)
	result4 := slice.Drop([]string{"a", "b", "c"}, 4)
	result5 := slice.Drop([]string{"a", "b", "c"}, -100)
	result6 := slice.Drop([]string{"a", "b", "c"}, 100)

	assert.Equal(t, []string{"a", "b", "c"}, result1)
	assert.Equal(t, []string{"b", "c"}, result2)
	assert.Equal(t, []string{"a", "b", "c"}, result3)
	assert.Equal(t, []string{}, result4)
	assert.Equal(t, []string{"a", "b", "c"}, result5)
	assert.Equal(t, []string{}, result6)
}

// TestDropRight Drop n elements from the end of a slice.
// func DropRight[T any](slice []T, n int) []T
func TestDropRight(t *testing.T) {
	result1 := slice.DropRight([]string{"a", "b", "c"}, 0)
	result2 := slice.DropRight([]string{"a", "b", "c"}, 1)
	result3 := slice.DropRight([]string{"a", "b", "c"}, -1)
	result4 := slice.DropRight([]string{"a", "b", "c"}, 4)
	result5 := slice.DropRight([]string{"a", "b", "c"}, -100)
	result6 := slice.DropRight([]string{"a", "b", "c"}, 100)

	assert.Equal(t, []string{"a", "b", "c"}, result1)
	assert.Equal(t, []string{"a", "b"}, result2)
	assert.Equal(t, []string{"a", "b", "c"}, result3)
	assert.Equal(t, []string{}, result4)
	assert.Equal(t, []string{"a", "b", "c"}, result5)
	assert.Equal(t, []string{}, result6)
}

// TestDropWhile Drop n elements from the start of a slice while predicate function returns true.
// func DropWhile[T any](slice []T, predicate func(item T) bool) []T
// 注意：如果条件不相等则循环 break
func TestDropWhile(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	result1 := slice.DropWhile(numbers, func(n int) bool {
		return true
	})
	result2 := slice.DropWhile(numbers, func(n int) bool {
		// [0:]
		return n == 0
	})
	result3 := slice.DropWhile(numbers, func(n int) bool {
		// [1:]
		return n == 1
	})
	result4 := slice.DropWhile(numbers, func(n int) bool {
		// [0:] 第一个元素不相等，则取不相等元素的索引及之后的数据
		return n == 2
	})
	result5 := slice.DropWhile(numbers, func(n int) bool {
		return n != 1
	})
	result6 := slice.DropWhile(numbers, func(n int) bool {
		return n != 2
	})
	result7 := slice.DropWhile(numbers, func(n int) bool {
		return n != 3
	})

	assert.Equal(t, []int{}, result1, "result1")
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result2, "result2")
	assert.Equal(t, []int{2, 3, 4, 5}, result3, "result3")
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result4, "result4")
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result5, "result5")
	assert.Equal(t, []int{2, 3, 4, 5}, result6, "result6")
	assert.Equal(t, []int{3, 4, 5}, result7, "result7")
}

// DropRightWhile Drop n elements from the end of a slice while predicate function returns true.
// func DropRightWhile[T any](slice []T, predicate func(item T) bool) []T
func TestDropRightWhile(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	result1 := slice.DropRightWhile(numbers, func(n int) bool {
		return n != 2
	})
	result2 := slice.DropRightWhile(numbers, func(n int) bool {
		return true
	})
	result3 := slice.DropRightWhile(numbers, func(n int) bool {
		return n == 0
	})

	assert.Equal(t, []int{1, 2}, result1, "result1")
	assert.Equal(t, []int{}, result2, "result2")
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result3, "result3")
}
