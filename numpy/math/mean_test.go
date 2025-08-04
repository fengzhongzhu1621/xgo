package math

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestMean(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Mean([]float32{2.3, 3.3, 4, 5.3})
	result2 := lo.Mean([]int32{2, 3, 4, 5})
	result3 := lo.Mean([]uint32{2, 3, 4, 5})
	result4 := lo.Mean([]uint32{})

	is.Equal(result1, float32(3.7250001))
	is.Equal(result2, int32(3))
	is.Equal(result3, uint32(3))
	is.Equal(result4, uint32(0))
}

func TestMeanBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.MeanBy([]float32{2.3, 3.3, 4, 5.3}, func(n float32) float32 { return n })
	result2 := lo.MeanBy([]int32{2, 3, 4, 5}, func(n int32) int32 { return n })
	result3 := lo.MeanBy([]uint32{2, 3, 4, 5}, func(n uint32) uint32 { return n })
	result4 := lo.MeanBy([]uint32{}, func(n uint32) uint32 { return n })

	is.Equal(result1, float32(3.7250001))
	is.Equal(result2, int32(3))
	is.Equal(result3, uint32(3))
	is.Equal(result4, uint32(0))
}
