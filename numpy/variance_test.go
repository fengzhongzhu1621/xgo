package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// 返回数字的方差。
// func Variance[T constraints.Float | constraints.Integer](numbers []T) float64
func TestVariance(t *testing.T) {
	result1 := mathutil.Variance([]int{1, 2, 3, 4, 5})
	result2 := mathutil.Variance([]float64{1.1, 2.2, 3.3, 4.4, 5.5})

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 2
	// 2.42
}

// Returns the standard deviation of numbers.
// func StdDev[T constraints.Float | constraints.Integer](numbers []T) float64
func TestStdDev(t *testing.T) {
	result1 := mathutil.TruncRound(mathutil.StdDev([]int{1, 2, 3, 4, 5}), 2)
	result2 := mathutil.TruncRound(mathutil.StdDev([]float64{1.1, 2.2, 3.3, 4.4, 5.5}), 2)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 1.41
	// 1.55
}
