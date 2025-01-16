package slice

import (
	"fmt"
	"testing"

	"github.com/araujo88/lambda-go/pkg/predicate"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestNth(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1, err1 := lo.Nth([]int{0, 1, 2, 3}, 2)
	result2, err2 := lo.Nth([]int{0, 1, 2, 3}, -2)
	result3, err3 := lo.Nth([]int{0, 1, 2, 3}, 42)
	result4, err4 := lo.Nth([]int{}, 0)
	result5, err5 := lo.Nth([]int{42}, 0)
	result6, err6 := lo.Nth([]int{42}, -1)

	is.Equal(result1, 2)
	is.Equal(err1, nil)
	is.Equal(result2, 2)
	is.Equal(err2, nil)
	is.Equal(result3, 0)
	is.Equal(err3, fmt.Errorf("nth: 42 out of slice bounds"))
	is.Equal(result4, 0)
	is.Equal(err4, fmt.Errorf("nth: 0 out of slice bounds"))
	is.Equal(result5, 42)
	is.Equal(err5, nil)
	is.Equal(result6, 42)
	is.Equal(err6, nil)
}

// TestFindBy Iterates over elements of slice, returning the first one that passes a truth test on predicate function.If return T is nil or zero value then no items matched the predicate func. In contrast to Find or FindLast, its return value no longer requires dereferencing.
// 遍历切片中的元素，返回第一个通过谓词函数真值测试的元素。如果返回值 T 是 nil 或零值，那么没有元素匹配该谓词函数。与 Find 或 FindLast 不同，其返回值不再需要解引用。
// func FindBy[T any](slice []T, predicate func(index int, item T) bool) (v T, ok bool)
func TestFindBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		index := 0
		result1, ok1 := lo.Find([]string{"a", "b", "c", "d"}, func(item string) bool {
			is.Equal([]string{"a", "b", "c", "d"}[index], item)
			index++
			return item == "b"
		})
		is.Equal(ok1, true)
		is.Equal(result1, "b")

		result2, ok2 := lo.Find([]string{"foobar"}, func(item string) bool {
			is.Equal("foobar", item)
			return item == "b"
		})
		is.Equal(ok2, false)
		is.Equal(result2, "")
	}

	{
		index := 0
		item1, index1, ok1 := lo.FindIndexOf([]string{"a", "b", "c", "d", "b"}, func(item string) bool {
			is.Equal([]string{"a", "b", "c", "d", "b"}[index], item)
			index++
			return item == "b"
		})
		is.Equal(item1, "b")
		is.Equal(ok1, true)
		is.Equal(index1, 1)

		item2, index2, ok2 := lo.FindIndexOf([]string{"foobar"}, func(item string) bool {
			is.Equal("foobar", item)
			return item == "b"
		})
		is.Equal(item2, "")
		is.Equal(ok2, false)
		is.Equal(index2, -1)
	}

	{
		nums := []int{1, 2, 3, 4, 5}

		isEven := func(i, num int) bool {
			return num%2 == 0
		}

		result, ok := slice.FindBy(nums, isEven)

		assert.Equal(t, 2, result)
		assert.Equal(t, true, ok)
	}

	{
		tests := []struct {
			name      string
			slice     []int
			predicate func(int) bool
			want      int
			found     bool
		}{
			{"finds element", []int{1, 2, 3}, func(x int) bool { return x == 3 }, 3, true},
			{"does not find element", []int{1, 2, 3}, func(x int) bool { return x == 5 }, 0, false},
			{"empty slice", []int{}, func(x int) bool { return true }, 0, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, found := predicate.Find(tt.slice, tt.predicate)
				if got != tt.want || found != tt.found {
					t.Errorf("Find() = %v, %v, want %v, %v", got, found, tt.want, tt.found)
				}
			})
		}
	}
}

// TestFindLastBy FindLastBy iterates over elements of slice, returning the last one that passes a truth test on predicate function. If return T is nil or zero value then no items matched the predicate func. In contrast to Find or FindLast, its return value no longer requires dereferencing.
// 遍历切片中的元素，返回第一个通过谓词函数真值测试的元素。如果返回值 T 是 nil 或零值，那么没有元素匹配该谓词函数。与 Find 或 FindLast 不同，其返回值不再需要解引用。
// func FindLastBy[T any](slice []T, predicate func(index int, item T) bool) (v T, ok bool)
func TestFindLastBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		index := 0
		item1, index1, ok1 := lo.FindLastIndexOf([]string{"a", "b", "c", "d", "b"}, func(item string) bool {
			is.Equal([]string{"b", "d", "c", "b", "a"}[index], item)
			index++
			return item == "b"
		})
		item2, index2, ok2 := lo.FindLastIndexOf([]string{"foobar"}, func(item string) bool {
			is.Equal("foobar", item)
			return item == "b"
		})

		is.Equal(item1, "b")
		is.Equal(ok1, true)
		is.Equal(index1, 4)
		is.Equal(item2, "")
		is.Equal(ok2, false)
		is.Equal(index2, -1)
	}

	{
		nums := []int{1, 2, 3, 4, 5}

		isEven := func(i, num int) bool {
			return num%2 == 0
		}

		result, ok := slice.FindLastBy(nums, isEven)

		assert.Equal(t, 4, result)
		assert.Equal(t, true, ok)
	}
}

// TestIndexOf Returns the index at which the first occurrence of a item is found in a slice or return -1 if the item cannot be found.
// func IndexOf[T comparable](slice []T, item T) int
func TestIndexOf(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.IndexOf([]int{0, 1, 2, 1, 2, 3}, 2)
		result2 := lo.IndexOf([]int{0, 1, 2, 1, 2, 3}, 6)

		is.Equal(result1, 2)
		is.Equal(result2, -1)
	}

	{
		strs := []string{"a", "a", "b", "c"}

		result1 := slice.IndexOf(strs, "a")
		result2 := slice.IndexOf(strs, "d")

		assert.Equal(t, 0, result1)
		assert.Equal(t, -1, result2)
	}
}

// TestLastIndexOf Returns the index at which the last occurrence of a item is found in a slice or return -1 if the item cannot be found.
// func LastIndexOf[T comparable](slice []T, item T) int
func TestLastIndexOf(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.LastIndexOf([]int{0, 1, 2, 1, 2, 3}, 2)
		result2 := lo.LastIndexOf([]int{0, 1, 2, 1, 2, 3}, 6)

		is.Equal(result1, 4)
		is.Equal(result2, -1)
	}
	{
		strs := []string{"a", "a", "b", "c"}

		result1 := slice.LastIndexOf(strs, "a")
		result2 := slice.LastIndexOf(strs, "d")

		assert.Equal(t, 1, result1)
		assert.Equal(t, -1, result2)
	}
}

func TestFindOrElse(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	index := 0
	result1 := lo.FindOrElse([]string{"a", "b", "c", "d"}, "x", func(item string) bool {
		is.Equal([]string{"a", "b", "c", "d"}[index], item)
		index++
		return item == "b"
	})
	result2 := lo.FindOrElse([]string{"foobar"}, "x", func(item string) bool {
		is.Equal("foobar", item)
		return item == "b"
	})

	is.Equal(result1, "b")
	is.Equal(result2, "x")
}
