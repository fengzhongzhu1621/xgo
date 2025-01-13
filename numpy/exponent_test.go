package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// TestExponent Calculate x to the nth power.
// func Exponent(x, n int64) int64
func TestExponent(t *testing.T) {
	result1 := mathutil.Exponent(10, 0)
	result2 := mathutil.Exponent(10, 1)
	result3 := mathutil.Exponent(10, 2)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 1
	// 10
	// 100
}
