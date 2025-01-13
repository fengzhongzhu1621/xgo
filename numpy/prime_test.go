package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// 检查数字是否为质数。
// func IsPrime(n int) bool
func TestIsPrime(t *testing.T) {
	result1 := mathutil.IsPrime(-1)
	result2 := mathutil.IsPrime(0)
	result3 := mathutil.IsPrime(1)
	result4 := mathutil.IsPrime(2)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)

	// Output:
	// false
	// false
	// false
	// true
}
