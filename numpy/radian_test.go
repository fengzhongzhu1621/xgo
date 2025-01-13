package numpy

import (
	"fmt"
	"math"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

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

// func RadianToAngle(radian float64) float64
func TestRadianToAngle(t *testing.T) {
	result1 := mathutil.RadianToAngle(math.Pi)
	result2 := mathutil.RadianToAngle(math.Pi / 2)
	result3 := mathutil.RadianToAngle(math.Pi / 4)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 180
	// 90
	// 45
}

// func Cos(radian float64, precision ...int) float64
func TestCos(t *testing.T) {
	result1 := mathutil.Cos(0)
	result2 := mathutil.Cos(90)
	result3 := mathutil.Cos(180)
	result4 := mathutil.Cos(math.Pi)
	result5 := mathutil.Cos(math.Pi / 2)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)
	fmt.Println(result5)

	// Output:
	// 1
	// -0.447
	// -0.598
	// -1
	// 0
}

// func Sin(radian float64, precision ...int) float64
func TestSin(t *testing.T) {
	result1 := mathutil.Sin(0)
	result2 := mathutil.Sin(90)
	result3 := mathutil.Sin(180)
	result4 := mathutil.Sin(math.Pi)
	result5 := mathutil.Sin(math.Pi / 2)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)
	fmt.Println(result5)

	// Output:
	// 0
	// 0.894
	// -0.801
	// 0
	// 1
}
