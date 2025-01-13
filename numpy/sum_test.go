package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// func Sum[T constraints.Integer | constraints.Float](numbers ...T) T
func TestSum(t *testing.T) {
	result1 := mathutil.Sum(1, 2)
	result2 := mathutil.Sum(0.1, float64(1))

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 3
	// 1.1
}
