package math

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// Calculates P(n, k).
// func Permutation(n, k uint) uint
func TestPermutation(t *testing.T) {
	result1 := mathutil.Permutation(5, 3)
	result2 := mathutil.Permutation(5, 5)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 60
	// 120
}

// Calculates C(n, k).
// func Combination(n, k uint) uint
func TestCombination(t *testing.T) {
	result1 := mathutil.Combination(5, 3)
	result2 := mathutil.Combination(5, 5)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 10
	// 1
}
