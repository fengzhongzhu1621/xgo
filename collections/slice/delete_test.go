package slice

import (
	"strconv"
	"strings"
	"testing"

	"github.com/samber/lo"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestDeleteAt delete the element of slice at index.
// func DeleteAt[T any](slice []T, index int)
func TestDeleteAt(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		is.Equal([]int{1, 2, 3, 4}, lo.DropByIndex([]int{0, 1, 2, 3, 4}, 0))
		is.Equal([]int{3, 4}, lo.DropByIndex([]int{0, 1, 2, 3, 4}, 0, 1, 2))
		is.Equal([]int{0, 4}, lo.DropByIndex([]int{0, 1, 2, 3, 4}, -4, -2, -3))
		is.Equal([]int{0, 2, 3, 4}, lo.DropByIndex([]int{0, 1, 2, 3, 4}, -4, -4))
		is.Equal([]int{2, 4}, lo.DropByIndex([]int{0, 1, 2, 3, 4}, 3, 1, 0))
		is.Equal([]int{0, 1, 3, 4}, lo.DropByIndex([]int{0, 1, 2, 3, 4}, 2))
		is.Equal([]int{0, 1, 2, 3}, lo.DropByIndex([]int{0, 1, 2, 3, 4}, 4))
		is.Equal([]int{0, 1, 2, 3, 4}, lo.DropByIndex([]int{0, 1, 2, 3, 4}, 5))
		is.Equal([]int{0, 1, 2, 3, 4}, lo.DropByIndex([]int{0, 1, 2, 3, 4}, 100))
		is.Equal([]int{0, 1, 2, 3}, lo.DropByIndex([]int{0, 1, 2, 3, 4}, -1))
		is.Equal([]int{}, lo.DropByIndex([]int{}, 0, 1))
		is.Equal([]int{}, lo.DropByIndex([]int{42}, 0, 1))
		is.Equal([]int{}, lo.DropByIndex([]int{42}, 1, 0))
		is.Equal([]int{}, lo.DropByIndex([]int{}, 1))
		is.Equal([]int{}, lo.DropByIndex([]int{1}, 0))
	}

	{
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
	t.Parallel()
	is := assert.New(t)

	{
		is.Equal([]int{1, 2, 3, 4}, lo.Drop([]int{0, 1, 2, 3, 4}, 1))
		is.Equal([]int{2, 3, 4}, lo.Drop([]int{0, 1, 2, 3, 4}, 2))
		is.Equal([]int{3, 4}, lo.Drop([]int{0, 1, 2, 3, 4}, 3))
		is.Equal([]int{4}, lo.Drop([]int{0, 1, 2, 3, 4}, 4))
		is.Equal([]int{}, lo.Drop([]int{0, 1, 2, 3, 4}, 5))
		is.Equal([]int{}, lo.Drop([]int{0, 1, 2, 3, 4}, 6))

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Drop(allStrings, 2)
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
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
}

// TestDropRight Drop n elements from the end of a slice.
// func DropRight[T any](slice []T, n int) []T
func TestDropRight(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		is.Equal([]int{0, 1, 2, 3}, lo.DropRight([]int{0, 1, 2, 3, 4}, 1))
		is.Equal([]int{0, 1, 2}, lo.DropRight([]int{0, 1, 2, 3, 4}, 2))
		is.Equal([]int{0, 1}, lo.DropRight([]int{0, 1, 2, 3, 4}, 3))
		is.Equal([]int{0}, lo.DropRight([]int{0, 1, 2, 3, 4}, 4))
		is.Equal([]int{}, lo.DropRight([]int{0, 1, 2, 3, 4}, 5))
		is.Equal([]int{}, lo.DropRight([]int{0, 1, 2, 3, 4}, 6))

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.DropRight(allStrings, 2)
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
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
}

// TestDropWhile Drop n elements from the start of a slice while predicate function returns true.
// func DropWhile[T any](slice []T, predicate func(item T) bool) []T
// 注意：如果条件不相等则循环 break
func TestDropWhile(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		is.Equal([]int{4, 5, 6}, lo.DropWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
			return t != 4
		}))

		is.Equal([]int{}, lo.DropWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
			return true
		}))

		is.Equal([]int{0, 1, 2, 3, 4, 5, 6}, lo.DropWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
			return t == 10
		}))

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.DropWhile(allStrings, func(t string) bool {
			return t != "foo"
		})
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
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
}

// DropRightWhile Drop n elements from the end of a slice while predicate function returns true.
// func DropRightWhile[T any](slice []T, predicate func(item T) bool) []T
func TestDropRightWhile(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		is.Equal([]int{0, 1, 2, 3}, lo.DropRightWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
			return t != 3
		}))

		is.Equal([]int{0, 1}, lo.DropRightWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
			return t != 1
		}))

		is.Equal([]int{0, 1, 2, 3, 4, 5, 6}, lo.DropRightWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
			return t == 10
		}))

		is.Equal([]int{}, lo.DropRightWhile([]int{0, 1, 2, 3, 4, 5, 6}, func(t int) bool {
			return t != 10
		}))

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.DropRightWhile(allStrings, func(t string) bool {
			return t != "foo"
		})
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
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
}

func TestReject(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.Reject([]int{1, 2, 3, 4}, func(x int, _ int) bool {
		return x%2 == 0
	})
	is.Equal(r1, []int{1, 3})

	r2 := lo.Reject([]string{"Smith", "foo", "Domin", "bar", "Olivia"}, func(x string, _ int) bool {
		return len(x) > 3
	})
	is.Equal(r2, []string{"foo", "bar"})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.Reject(allStrings, func(x string, _ int) bool {
		return len(x) > 0
	})
	is.IsType(nonempty, allStrings, "type preserved")
}

func TestRejectMap(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.RejectMap([]int64{1, 2, 3, 4}, func(x int64, _ int) (string, bool) {
		if x%2 == 0 {
			return strconv.FormatInt(x, 10), false
		}
		return "", true
	})
	r2 := lo.RejectMap([]string{"cpu", "gpu", "mouse", "keyboard"}, func(x string, _ int) (string, bool) {
		if strings.HasSuffix(x, "pu") {
			return "xpu", false
		}
		return "", true
	})

	is.Equal(len(r1), 2)
	is.Equal(len(r2), 2)
	is.Equal(r1, []string{"2", "4"})
	is.Equal(r2, []string{"xpu", "xpu"})
}

func TestFilterReject(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	left1, right1 := lo.FilterReject([]int{1, 2, 3, 4}, func(x int, _ int) bool {
		return x%2 == 0
	})

	is.Equal(left1, []int{2, 4})
	is.Equal(right1, []int{1, 3})

	left2, right2 := lo.FilterReject([]string{"Smith", "foo", "Domin", "bar", "Olivia"}, func(x string, _ int) bool {
		return len(x) > 3
	})

	is.Equal(left2, []string{"Smith", "Domin", "Olivia"})
	is.Equal(right2, []string{"foo", "bar"})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	a, b := lo.FilterReject(allStrings, func(x string, _ int) bool {
		return len(x) > 0
	})
	is.IsType(a, allStrings, "type preserved")
	is.IsType(b, allStrings, "type preserved")
}
