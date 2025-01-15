package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestIntersection Creates a slice of unique values that included by all slices.
// func Intersection[T comparable](slices ...[]T) []T
func TestIntersection(t *testing.T) {
	nums1 := []int{1, 2, 3}
	nums2 := []int{2, 3, 4}

	result := slice.Intersection(nums1, nums2)

	assert.Equal(t, []int{2, 3}, result)
}

func TestIntersect2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
	result2 := lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
	result3 := lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
	result4 := lo.Intersect([]int{0, 6}, []int{0, 1, 2, 3, 4, 5})
	result5 := lo.Intersect([]int{0, 6, 0}, []int{0, 1, 2, 3, 4, 5})

	is.Equal(result1, []int{0, 2})
	is.Equal(result2, []int{0})
	is.Equal(result3, []int{})
	is.Equal(result4, []int{0})
	is.Equal(result5, []int{0})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.Intersect(allStrings, allStrings)
	is.IsType(nonempty, allStrings, "type preserved")
}

// TestUnion Creates a slice of unique values, in order, from all given slices. using == for equality comparisons.
// func Union[T comparable](slices ...[]T) []T
func TestUnion(t *testing.T) {
	nums1 := []int{1, 3, 4, 6}
	nums2 := []int{1, 2, 5, 6}

	result := slice.Union(nums1, nums2)

	assert.Equal(t, []int{1, 3, 4, 6, 2, 5}, result)
}

// TestUnionBy UnionBy is like Union, what's more it accepts iteratee which is invoked for each element of each slice.
// func UnionBy[T any, V comparable](predicate func(item T) V, slices ...[]T) []T
func TestUnionBy(t *testing.T) {
	nums := []int{1, 2, 3, 4}

	divideTwo := func(n int) int {
		return n / 2
	}
	result := slice.UnionBy(divideTwo, nums)

	assert.Equal(t, []int{1, 2, 4}, result)
}

func TestUnion2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Union([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 10})
	result2 := lo.Union([]int{0, 1, 2, 3, 4, 5}, []int{6, 7})
	result3 := lo.Union([]int{0, 1, 2, 3, 4, 5}, []int{})
	result4 := lo.Union([]int{0, 1, 2}, []int{0, 1, 2})
	result5 := lo.Union([]int{}, []int{})
	is.Equal(result1, []int{0, 1, 2, 3, 4, 5, 10})
	is.Equal(result2, []int{0, 1, 2, 3, 4, 5, 6, 7})
	is.Equal(result3, []int{0, 1, 2, 3, 4, 5})
	is.Equal(result4, []int{0, 1, 2})
	is.Equal(result5, []int{})

	result11 := lo.Union([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 10}, []int{0, 1, 11})
	result12 := lo.Union([]int{0, 1, 2, 3, 4, 5}, []int{6, 7}, []int{8, 9})
	result13 := lo.Union([]int{0, 1, 2, 3, 4, 5}, []int{}, []int{})
	result14 := lo.Union([]int{0, 1, 2}, []int{0, 1, 2}, []int{0, 1, 2})
	result15 := lo.Union([]int{}, []int{}, []int{})
	is.Equal(result11, []int{0, 1, 2, 3, 4, 5, 10, 11})
	is.Equal(result12, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	is.Equal(result13, []int{0, 1, 2, 3, 4, 5})
	is.Equal(result14, []int{0, 1, 2})
	is.Equal(result15, []int{})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.Union(allStrings, allStrings)
	is.IsType(nonempty, allStrings, "type preserved")
}

// TestDifference Creates an slice of whose element not included in the other given slice.
// func Difference[T comparable](slice, comparedSlice []T) []T
func TestDifference(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{4, 5, 6}

	result := slice.Difference(s1, s2)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestDifference2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	left1, right1 := lo.Difference([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 6})
	is.Equal(left1, []int{1, 3, 4, 5})
	is.Equal(right1, []int{6})

	left2, right2 := lo.Difference([]int{1, 2, 3, 4, 5}, []int{0, 6})
	is.Equal(left2, []int{1, 2, 3, 4, 5})
	is.Equal(right2, []int{0, 6})

	left3, right3 := lo.Difference([]int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5})
	is.Equal(left3, []int{})
	is.Equal(right3, []int{})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	a, b := lo.Difference(allStrings, allStrings)
	is.IsType(a, allStrings, "type preserved")
	is.IsType(b, allStrings, "type preserved")
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

// TestSymmetricDifference Create a slice whose element is in given slices, but not in both slices.
// func SymmetricDifference[T comparable](slices ...[]T) []T
func TestSymmetricDifference(t *testing.T) {
	nums1 := []int{1, 2, 3}
	nums2 := []int{1, 2, 4}

	result := slice.SymmetricDifference(nums1, nums2)

	assert.Equal(t, []int{3, 4}, result)
}
