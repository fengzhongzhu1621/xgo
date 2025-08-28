package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

type foo struct {
	bar string
}

func (f foo) Clone() foo {
	return foo{f.bar}
}

// fills elements of array with `initial` value.
func TestFill(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Fill([]foo{{"a"}, {"a"}}, foo{"b"})
	result2 := lo.Fill([]foo{}, foo{"a"})

	is.Equal(result1, []foo{{"b"}, {"b"}})
	is.Equal(result2, []foo{})
}

// TestReplace Returns a copy of the slice with the first n non-overlapping instances of old replaced by new.
// 返回一个切片副本，其中前n个不重叠的旧实例被新实例替换。
// func Replace[T comparable](slice []T, old T, new T, n int) []T
// n 表示被替换的数量
func TestReplace(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		in := []int{0, 1, 0, 1, 2, 3, 0}

		out1 := lo.Replace(in, 0, 42, 2)
		out2 := lo.Replace(in, 0, 42, 1)
		out3 := lo.Replace(in, 0, 42, 0)
		out4 := lo.Replace(in, 0, 42, -1)
		out5 := lo.Replace(in, 0, 42, -1)
		out6 := lo.Replace(in, -1, 42, 2)
		out7 := lo.Replace(in, -1, 42, 1)
		out8 := lo.Replace(in, -1, 42, 0)
		out9 := lo.Replace(in, -1, 42, -1)
		out10 := lo.Replace(in, -1, 42, -1)

		is.Equal([]int{42, 1, 42, 1, 2, 3, 0}, out1)
		is.Equal([]int{42, 1, 0, 1, 2, 3, 0}, out2)
		is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, out3)
		is.Equal([]int{42, 1, 42, 1, 2, 3, 42}, out4)
		is.Equal([]int{42, 1, 42, 1, 2, 3, 42}, out5)
		is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, out6)
		is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, out7)
		is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, out8)
		is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, out9)
		is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, out10)

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Replace(allStrings, "0", "2", 1)
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
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
}

// TestReplaceAll Returns a copy of the slice with the first n non-overlapping instances of old replaced by new.
// func ReplaceAll[T comparable](slice []T, old T, new T) []T
func TestReplaceAll(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		in := []int{0, 1, 0, 1, 2, 3, 0}

		out1 := lo.ReplaceAll(in, 0, 42)
		out2 := lo.ReplaceAll(in, -1, 42)

		is.Equal([]int{42, 1, 42, 1, 2, 3, 42}, out1)
		is.Equal([]int{0, 1, 0, 1, 2, 3, 0}, out2)

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.ReplaceAll(allStrings, "0", "2")
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
		result := slice.ReplaceAll([]string{"a", "b", "c", "a"}, "a", "x")
		assert.Equal(t, []string{"x", "b", "c", "x"}, result)
	}
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
	modifiedStrs, count := slice.SetToDefaultIf(strs, func(s string) bool { return s == "a" })

	assert.Equal(t, []string{"", "b", "", "c", "d", ""}, modifiedStrs)
	assert.Equal(t, 3, count)
}
