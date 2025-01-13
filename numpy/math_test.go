package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// TestAverage Return average value of numbers. Maybe call RoundToFloat to round result.
// func Average[T constraints.Integer | constraints.Float](numbers ...T) float64
// func RoundToFloat[T constraints.Float | constraints.Integer](x T, n int) float64
func TestAverage(t *testing.T) {
	result1 := mathutil.Average(1, 2)

	avg := mathutil.Average(1.2, 1.25)
	result2 := mathutil.RoundToFloat(avg, 1)
	result3 := mathutil.RoundToFloat(avg, 2)

	fmt.Println(result1) // 1.5
	fmt.Println(result2) // 1.2
	fmt.Println(result3) // 1.23
}

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

func TestFibonacci(t *testing.T) {
	result1 := mathutil.Fibonacci(1, 1, 1)
	result2 := mathutil.Fibonacci(1, 1, 2)
	result3 := mathutil.Fibonacci(1, 1, 5)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 1
	// 1
	// 5
}

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

// func Max[T constraints.Integer | constraints.Float](numbers ...T) T
func TestMax(t *testing.T) {
	result1 := mathutil.Max(1, 2, 3)
	result2 := mathutil.Max(1.2, 1.4, 1.1, 1.4)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 3
	// 1.4
}

// func MaxBy[T any](slice []T, comparator func(T, T) bool) T
func TestMaxBy(t *testing.T) {
	result1 := mathutil.MaxBy([]string{"a", "ab", "abc"}, func(v1, v2 string) bool {
		return len(v1) > len(v2)
	})

	result2 := mathutil.MaxBy([]string{"abd", "abc", "ab"}, func(v1, v2 string) bool {
		return len(v1) > len(v2)
	})

	result3 := mathutil.MaxBy([]string{}, func(v1, v2 string) bool {
		return len(v1) > len(v2)
	})

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// abc
	// abd
	//
}

// func Min[T constraints.Integer | constraints.Float](numbers ...T) T
func TestMin(t *testing.T) {
	result1 := mathutil.Min(1, 2, 3)
	result2 := mathutil.Min(1.2, 1.4, 1.1, 1.4)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 1
	// 1.1
}

// func MinBy[T any](slice []T, comparator func(T, T) bool) T
func TestMinBy(t *testing.T) {
	result1 := mathutil.MinBy([]string{"a", "ab", "abc"}, func(v1, v2 string) bool {
		return len(v1) < len(v2)
	})

	result2 := mathutil.MinBy([]string{"ab", "ac", "abc"}, func(v1, v2 string) bool {
		return len(v1) < len(v2)
	})

	result3 := mathutil.MinBy([]string{}, func(v1, v2 string) bool {
		return len(v1) < len(v2)
	})

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// a
	// ab
	//
}

// TestPercent 计算val占总数的百分比，保留n位小数。
// func Percent(val, total float64, n int) float64
func TestPercent(t *testing.T) {
	result1 := mathutil.Percent(1, 2, 2)
	result2 := mathutil.Percent(0.1, 0.3, 2)
	result3 := mathutil.Percent(-30305, 408420, 2)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 50
	// 33.33
	// -7.42
}

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

// TestRange 从起始值start开始创建一个由指定数量count的数字组成的切片，元素步长为1。
// func Range[T constraints.Integer | constraints.Float](start T, count int) []T
func TestRange(t *testing.T) {
	result1 := mathutil.Range(1, 4)
	result2 := mathutil.Range(1, -4)
	result3 := mathutil.Range(-4, 4)
	result4 := mathutil.Range(1.0, 4)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)

	// Output:
	// [1 2 3 4]
	// [1 2 3 4]
	// [-4 -3 -2 -1]
	// [1 2 3 4]
}

// func RangeWithStep[T constraints.Integer | constraints.Float](start, end, step T) []T
func TestRangeWithStep(t *testing.T) {
	result1 := mathutil.RangeWithStep(1, 4, 1)
	result2 := mathutil.RangeWithStep(1, -1, 0)
	result3 := mathutil.RangeWithStep(-4, 1, 2)
	result4 := mathutil.RangeWithStep(1.0, 4.0, 1.1)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)

	// Output:
	// [1 2 3]
	// []
	// [-4 -2 0]
	// [1 2.1 3.2]
}

// TestAngleToRadian 将角度值转换为弧度值。
// func AngleToRadian(angle float64) float64
func TestAngleToRadian(t *testing.T) {
	result1 := mathutil.AngleToRadian(45)
	result2 := mathutil.AngleToRadian(90)
	result3 := mathutil.AngleToRadian(180)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 0.7853981633974483
	// 1.5707963267948966
	// 3.141592653589793
}
