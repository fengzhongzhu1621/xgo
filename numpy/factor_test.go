package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// TestFactorial 计算x的阶乘。
func TestFactorial(t *testing.T) {
	result1 := mathutil.Factorial(1)
	result2 := mathutil.Factorial(2)
	result3 := mathutil.Factorial(3)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 1
	// 2
	// 6
}
