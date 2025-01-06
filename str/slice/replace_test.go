package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestReplace Returns a copy of the slice with the first n non-overlapping instances of old replaced by new.
// 返回一个切片副本，其中前n个不重叠的旧实例被新实例替换。
// func Replace[T comparable](slice []T, old T, new T, n int) []T
// n 表示被替换的数量
func TestReplace(t *testing.T) {
	strs := []string{"a", "b", "c", "a"}

	result1 := slice.Replace(strs, "a", "x", 0)
	result2 := slice.Replace(strs, "a", "x", 1)
	result3 := slice.Replace(strs, "a", "x", 2)
	result4 := slice.Replace(strs, "a", "x", 3)
	result5 := slice.Replace(strs, "a", "x", -1)

	assert.Equal(t, []string{"a", "b", "c", "a"}, result1)
	assert.Equal(t, []string{"x", "b", "c", "a"}, result2)
	assert.Equal(t, []string{"x", "b", "c", "x"}, result3)
	assert.Equal(t, []string{"x", "b", "c", "x"}, result4)
	assert.Equal(t, []string{"x", "b", "c", "x"}, result5)
}

// TestReplaceAll Returns a copy of the slice with the first n non-overlapping instances of old replaced by new.
// func ReplaceAll[T comparable](slice []T, old T, new T) []T
func TestReplaceAll(t *testing.T) {
	result := slice.ReplaceAll([]string{"a", "b", "c", "a"}, "a", "x")

	assert.Equal(t, []string{"x", "b", "c", "x"}, result)
}

// TestUpdateAt Update the slice element at index. if param index < 0 or index <= len(slice), will return error.
// func UpdateAt[T any](slice []T, index int, value T) []T
func TestUpdateAt(t *testing.T) {
	result1 := slice.UpdateAt([]string{"a", "b", "c"}, -1, "1")
	result2 := slice.UpdateAt([]string{"a", "b", "c"}, 0, "1")
	result3 := slice.UpdateAt([]string{"a", "b", "c"}, 1, "1")
	result4 := slice.UpdateAt([]string{"a", "b", "c"}, 2, "1")
	result5 := slice.UpdateAt([]string{"a", "b", "c"}, 3, "1")

	assert.Equal(t, []string{"a", "b", "c"}, result1)
	assert.Equal(t, []string{"1", "b", "c"}, result2)
	assert.Equal(t, []string{"a", "1", "c"}, result3)
	assert.Equal(t, []string{"a", "b", "1"}, result4)
	assert.Equal(t, []string{"a", "b", "c"}, result5)
}

// TestSetToDefaultIf Sets elements to their default value if they match the given predicate. It retains the positions of the elements in the slice. It returns slice of T and the count of modified slice items
// 将匹配的值设置为默认值
// func SetToDefaultIf[T any](slice []T, predicate func(T) bool) ([]T, int)
func TestSetToDefaultIf(t *testing.T) {
	strs := []string{"a", "b", "a", "c", "d", "a"}
	modifiedStrs, count := slice.SetToDefaultIf(strs, func(s string) bool { return "a" == s })

	assert.Equal(t, []string{"", "b", "", "c", "d", ""}, modifiedStrs)
	assert.Equal(t, 3, count)
}
