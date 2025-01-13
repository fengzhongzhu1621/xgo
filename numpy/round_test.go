package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// func RoundToFloat[T constraints.Float | constraints.Integer](x T, n int) float64
func TestRoundToFloat(t *testing.T) {
	result1 := mathutil.RoundToFloat(0.124, 2)
	result2 := mathutil.RoundToFloat(0.125, 2)
	result3 := mathutil.RoundToFloat(0.125, 3)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 0.12
	// 0.13
	// 0.125
}

// func RoundToString[T constraints.Float | constraints.Integer](x T, n int) string
func TestRoundToString(t *testing.T) {
	result1 := mathutil.RoundToString(0.124, 2)
	result2 := mathutil.RoundToString(0.125, 2)
	result3 := mathutil.RoundToString(0.125, 3)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 0.12
	// 0.13
	// 0.125
}

// func TruncRound[T constraints.Float | constraints.Integer](x T, n int) T
func TestTruncRound(t *testing.T) {
	result1 := mathutil.TruncRound(0.124, 2)
	result2 := mathutil.TruncRound(0.125, 2)
	result3 := mathutil.TruncRound(0.125, 3)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 0.12
	// 0.12
	// 0.125
}
