package math

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// func Abs[T constraints.Integer | constraints.Float](x T) T
func TestAbs(t *testing.T) {
	result1 := mathutil.Abs(-1)
	result2 := mathutil.Abs(-0.1)
	result3 := mathutil.Abs(float32(0.2))

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 1
	// 0.1
	// 0.2
}
