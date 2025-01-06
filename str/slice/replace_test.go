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
