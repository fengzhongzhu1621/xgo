package slice

import (
	"math/rand"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestContain 判断 slice 是否包含指定值
// func Contain[T comparable](slice []T, target T) bool
func TestContain(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		rand.Seed(time.Now().UnixNano())

		result1 := lo.Sample([]string{"a", "b", "c"})
		result2 := lo.Sample([]string{})

		is.True(lo.Contains([]string{"a", "b", "c"}, result1))
		is.Equal(result2, "")
	}

	{
		result1 := slice.Contain([]string{"a", "b", "c"}, "a")
		result2 := slice.Contain([]int{1, 2, 3}, 4)

		assert.Equal(t, true, result1)
		assert.Equal(t, false, result2)
	}
}

// TestContainBy returns true if predicate function return true.
// func ContainBy[T any](slice []T, predicate func(item T) bool) bool
func TestContainBy(t *testing.T) {
	type foo struct {
		A string
		B int
	}

	array1 := []foo{{A: "1", B: 1}, {A: "2", B: 2}}
	result1 := slice.ContainBy(array1, func(f foo) bool { return f.A == "1" && f.B == 1 })
	result2 := slice.ContainBy(array1, func(f foo) bool { return f.A == "2" && f.B == 1 })

	array2 := []string{"a", "b", "c"}
	result3 := slice.ContainBy(array2, func(t string) bool { return t == "a" })
	result4 := slice.ContainBy(array2, func(t string) bool { return t == "d" })

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, true, result3)
	assert.Equal(t, false, result4)
}

// TestContainSubSlice Check if the slice contain subslice or not.
// func ContainSubSlice[T comparable](slice, subSlice []T) bool
func TestContainSubSlice(t *testing.T) {
	result1 := slice.ContainSubSlice([]string{"a", "b", "c"}, []string{"a", "b"})
	result2 := slice.ContainSubSlice([]string{"a", "b", "c"}, []string{"a", "d"})

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}
