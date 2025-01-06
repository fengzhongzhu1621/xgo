package randutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/slice"
)

// 生成随机int, 范围[min, max)
func TestRandomRandInt(t *testing.T) {
	rInt := random.RandInt(1, 10)
	fmt.Println(rInt)
}

func TestRandomRandString(t *testing.T) {
	randStr := random.RandString(6)
	fmt.Println(randStr)
}

// TestSliceRandom Random get a random item of slice, return idx=-1 when slice is empty.
// func Random[T any](slice []T) (val T, idx int)
func TestSliceRandom(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	val, idx := slice.Random(nums)
	if idx >= 0 && idx < len(nums) && slice.Contain(nums, val) {
		fmt.Println("okk")
	}
}
