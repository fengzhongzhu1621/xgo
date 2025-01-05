package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestDifference Creates an slice of whose element not included in the other given slice.
// func Difference[T comparable](slice, comparedSlice []T) []T
func TestDifference(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{4, 5, 6}

	result := slice.Difference(s1, s2)
	assert.Equal(t, []int{1, 2, 3}, result)
}

// TestDifferenceBy DifferenceBy accepts iteratee func which is invoked for each element of slice and values to generate the criterion by which they're compared.
// 接受一个迭代函数（iteratee func），该函数会针对切片中的每个元素被调用，并且接受一些值来生成用于比较它们的标准。
// func DifferenceBy[T comparable](slice []T, comparedSlice []T, iteratee func(index int, item T) T) []T
func TestDifferenceBy(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{3, 4, 5}

	addOne := func(index int, item int) int {
		// 对 s1 和 s2 切片中的元素都加上 1
		return item + 1
	}
	result := slice.DifferenceBy(s1, s2, addOne)
	assert.Equal(t, []int{1, 2}, result)
}

// TestDifferenceWith DifferenceWith accepts comparator which is invoked to compare elements of slice to values. The order and references of result values are determined by the first slice.
// func DifferenceWith[T any](slice []T, comparedSlice []T, comparator func(value, otherValue T) bool) []T
func TestDifferenceWith(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{4, 5, 6, 7, 8}

	isDouble := func(v1, v2 int) bool {
		return v2 == 2*v1
	}

	result := slice.DifferenceWith(s1, s2, isDouble)

	assert.Equal(t, []int{1, 5}, result)
}
