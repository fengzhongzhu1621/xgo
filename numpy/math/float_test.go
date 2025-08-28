package math

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// func CeilToFloat[T constraints.Float | constraints.Integer](x T, n int) float64
func TestCeilToFloat(t *testing.T) {
	result1 := mathutil.CeilToFloat(3.14159, 1)
	result2 := mathutil.CeilToFloat(3.14159, 2)
	result3 := mathutil.CeilToFloat(5, 4)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 3.2
	// 3.15
	// 5
}

// func CeilToString[T constraints.Float | constraints.Integer](x T, n int) string
func TestCeilToString(t *testing.T) {
	result1 := mathutil.CeilToString(3.14159, 1)
	result2 := mathutil.CeilToString(3.14159, 2)
	result3 := mathutil.CeilToString(5, 4)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 3.2
	// 3.15
	// 5.0000
}

// func FloorToFloat[T constraints.Float | constraints.Integer](x T, n int) float64
func TestFloorToFloat(t *testing.T) {
	result1 := mathutil.FloorToFloat(3.14159, 1)
	result2 := mathutil.FloorToFloat(3.14159, 2)
	result3 := mathutil.FloorToFloat(5, 4)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 3.1
	// 3.14
	// 5
}

// func FloorToString[T constraints.Float | constraints.Integer](x T, n int) string
func TestFloorToString(t *testing.T) {
	result1 := mathutil.FloorToString(3.14159, 1)
	result2 := mathutil.FloorToString(3.14159, 2)
	result3 := mathutil.FloorToString(5, 4)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 3.1
	// 3.14
	// 5.0000
}
