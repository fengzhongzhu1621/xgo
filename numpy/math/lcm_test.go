package math

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// Return Least Common Multiple (LCM) of integers.
// func LCM[T constraints.Integer](integers ...T) T
func TestLCM(t *testing.T) {
	result1 := mathutil.LCM(1, 1)
	result2 := mathutil.LCM(1, 2)
	result3 := mathutil.LCM(3, 6, 9)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 1
	// 2
	// 18
}
