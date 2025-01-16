package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// func Sum[T constraints.Integer | constraints.Float](numbers ...T) T
func TestSum(t *testing.T) {
	is := assert.New(t)

	{
		result1 := lo.Sum([]float32{2.3, 3.3, 4, 5.3})
		result2 := lo.Sum([]int32{2, 3, 4, 5})
		result3 := lo.Sum([]uint32{2, 3, 4, 5})
		result4 := lo.Sum([]uint32{})
		result5 := lo.Sum([]complex128{4_4, 2_2})

		is.Equal(result1, float32(14.900001))
		is.Equal(result2, int32(14))
		is.Equal(result3, uint32(14))
		is.Equal(result4, uint32(0))
		is.Equal(result5, complex128(6_6))
	}

	{
		result1 := mathutil.Sum(1, 2)
		result2 := mathutil.Sum(0.1, float64(1))

		fmt.Println(result1)
		fmt.Println(result2)

		// Output:
		// 3
		// 1.1
	}
}

func TestSumBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.SumBy([]float32{2.3, 3.3, 4, 5.3}, func(n float32) float32 { return n })
	result2 := lo.SumBy([]int32{2, 3, 4, 5}, func(n int32) int32 { return n })
	result3 := lo.SumBy([]uint32{2, 3, 4, 5}, func(n uint32) uint32 { return n })
	result4 := lo.SumBy([]uint32{}, func(n uint32) uint32 { return n })
	result5 := lo.SumBy([]complex128{4_4, 2_2}, func(n complex128) complex128 { return n })

	is.Equal(result1, float32(14.900001))
	is.Equal(result2, int32(14))
	is.Equal(result3, uint32(14))
	is.Equal(result4, uint32(0))
	is.Equal(result5, complex128(6_6))
}
