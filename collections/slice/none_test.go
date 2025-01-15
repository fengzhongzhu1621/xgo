package slice

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestNone(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.None([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
	result2 := lo.None([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
	result3 := lo.None([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
	result4 := lo.None([]int{0, 1, 2, 3, 4, 5}, []int{})

	is.False(result1)
	is.False(result2)
	is.True(result3)
	is.True(result4)
}

func TestNoneBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.NoneBy([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 5
	})

	is.False(result1)

	result2 := lo.NoneBy([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 3
	})

	is.False(result2)

	result3 := lo.NoneBy([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 0
	})

	is.True(result3)

	result4 := lo.NoneBy([]int{}, func(x int) bool {
		return x < 5
	})

	is.True(result4)
}
