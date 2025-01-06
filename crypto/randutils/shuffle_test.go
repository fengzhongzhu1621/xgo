package randutils

import (
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"testing"
)

// TestShuffle 打乱给定字符串的字符顺序。
// func Shuffle(str string) string
func TestShuffle(t *testing.T) {
	result1 := strutil.Shuffle("hello")
	result2 := strutil.Shuffle("hello")

	fmt.Println(result1)
	fmt.Println(result2)
}

// TestLancetShuffle Creates an slice of shuffled values.
// func Shuffle[T any](slice []T) []T
func TestLancetShuffle(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	result := slice.Shuffle(nums)

	fmt.Println(result)
}
