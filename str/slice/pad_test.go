package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestRightPadding adds padding to the right end of a slice.
// func RightPadding[T any](slice []T, paddingValue T, paddingLength int) []T
func TestRightPadding(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	padded := slice.RightPadding(nums, 0, 3)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 0, 0, 0}, padded)
}

// TestLeftPadding adds padding to the left begin of a slice.
// func LeftPadding[T any](slice []T, paddingValue T, paddingLength int) []T
func TestLeftPadding(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	padded := slice.LeftPadding(nums, 0, 3)

	assert.Equal(t, []int{0, 0, 0, 1, 2, 3, 4, 5}, padded)
}
