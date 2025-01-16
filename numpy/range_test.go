package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestRange 从起始值start开始创建一个由指定数量count的数字组成的切片，元素步长为1。
// func Range[T constraints.Integer | constraints.Float](start T, count int) []T
func TestRange(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.Range(4)
		result2 := lo.Range(-4)
		result3 := lo.Range(0)
		is.Equal(result1, []int{0, 1, 2, 3})
		is.Equal(result2, []int{0, -1, -2, -3})
		is.Equal(result3, []int{})
	}

	{
		result1 := lo.RangeFrom(1, 5)
		result2 := lo.RangeFrom(-1, -5)
		result3 := lo.RangeFrom(10, 0)
		result4 := lo.RangeFrom(2.0, 3)
		result5 := lo.RangeFrom(-2.0, -3)
		is.Equal(result1, []int{1, 2, 3, 4, 5})
		is.Equal(result2, []int{-1, -2, -3, -4, -5})
		is.Equal(result3, []int{})
		is.Equal(result4, []float64{2.0, 3.0, 4.0})
		is.Equal(result5, []float64{-2.0, -3.0, -4.0})
	}

	{
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
}

// func RangeWithStep[T constraints.Integer | constraints.Float](start, end, step T) []T
func TestRangeWithStep(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.RangeWithSteps(0, 20, 6)
		result2 := lo.RangeWithSteps(0, 3, -5)
		result3 := lo.RangeWithSteps(1, 1, 0)
		result4 := lo.RangeWithSteps(3, 2, 1)
		result5 := lo.RangeWithSteps(1.0, 4.0, 2.0)
		result6 := lo.RangeWithSteps[float32](-1.0, -4.0, -1.0)
		is.Equal([]int{0, 6, 12, 18}, result1)
		is.Equal([]int{}, result2)
		is.Equal([]int{}, result3)
		is.Equal([]int{}, result4)
		is.Equal([]float64{1.0, 3.0}, result5)
		is.Equal([]float32{-1.0, -2.0, -3.0}, result6)
	}

	{
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
}

// 将一个数值限制在给定的最小值和最大值之间。具体来说，lo.Clamp 函数接受三个参数：
// value：要限制的数值。
// min：允许的最小值。
// max：允许的最大值。
// 函数的返回值将是 value，但如果 value 小于 min，则返回 min；如果 value 大于 max，则返回 max。这样，返回的值将始终在 min 和 max 之间（包括 min 和 max）。
func TestClamp(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Clamp(0, -10, 10)
	result2 := lo.Clamp(-42, -10, 10)
	result3 := lo.Clamp(42, -10, 10)

	is.Equal(result1, 0)
	is.Equal(result2, -10)
	is.Equal(result3, 10)
}
