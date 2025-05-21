package math

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// 返回整数的最大公约数（GCD）。
// func GCD[T constraints.Integer](integers ...T) T
func TestGCD(t *testing.T) {
	result1 := mathutil.GCD(1, 1)
	result2 := mathutil.GCD(1, -1)
	result3 := mathutil.GCD(-1, 1)
	result4 := mathutil.GCD(-1, -1)
	result5 := mathutil.GCD(3, 6, 9)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)
	fmt.Println(result5)

	// Output:
	// 1
	// 1
	// -1
	// -1
	// 3
}
