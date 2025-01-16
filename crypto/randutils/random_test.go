package randutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/slice"
)

// TestSliceRandom Random get a random item of slice, return idx=-1 when slice is empty.
// func Random[T any](slice []T) (val T, idx int)
func TestSliceRandom(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	val, idx := slice.Random(nums)
	if idx >= 0 && idx < len(nums) && slice.Contain(nums, val) {
		fmt.Println("okk")
	}
}

// TestRandFromGivenSlice Generate a random element from given slice.
// func RandFromGivenSlice[T any](slice []T) T
func TestRandFromGivenSlice(t *testing.T) {
	randomSet := []any{"a", 8, "hello", true, 1.1}
	randElm := random.RandFromGivenSlice(randomSet)
	fmt.Println(randElm)
}
