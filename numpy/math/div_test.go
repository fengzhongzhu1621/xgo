package math

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// func Div[T constraints.Float | constraints.Integer](x T, y T) float64
func TestDiv(t *testing.T) {
	result1 := mathutil.Div(9, 4)
	result2 := mathutil.Div(1, 2)
	result3 := mathutil.Div(0, 666)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	// Output:
	// 2.25
	// 0.5
	// 0
}
