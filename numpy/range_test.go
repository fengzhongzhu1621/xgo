package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

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
