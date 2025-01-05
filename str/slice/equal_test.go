package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestEqual 检查两个切片是否相等，相等条件：切片长度相同，元素顺序和值都相同。
// func Equal[T comparable](slice1, slice2 []T) bool
func TestEqual(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{1, 3, 2}

	result1 := slice.Equal(s1, s2)
	result2 := slice.Equal(s1, s3)

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}

// TestEqualWith Check if two slices are equal with comparator func.
// func EqualWith[T, U any](slice1 []T, slice2 []U, comparator func(T, U) bool) bool
func TestEqualWith(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{2, 4, 6}

	isDouble := func(a, b int) bool {
		return b == a*2
	}

	result := slice.EqualWith(s1, s2, isDouble)

	assert.Equal(t, true, result)
}
