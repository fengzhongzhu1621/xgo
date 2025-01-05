package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
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
