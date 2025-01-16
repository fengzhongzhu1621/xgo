package slice

import (
	"math"
	"testing"

	"github.com/samber/lo"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestRepeat Creates a slice with length n whose elements are passed param item.
// func Repeat[T any](item T, n int) []T
func TestRepeat(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.Repeat(2, foo{"a"})
		result2 := lo.Repeat(0, foo{"a"})

		is.Equal(result1, []foo{{"a"}, {"a"}})
		is.Equal(result2, []foo{})
	}

	{
		result := slice.Repeat("a", 3)
		assert.Equal(t, []string{"a", "a", "a"}, result)
	}
}

func TestRepeatBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cb := func(i int) int {
		return int(math.Pow(float64(i), 2))
	}

	result1 := lo.RepeatBy(0, cb)
	result2 := lo.RepeatBy(2, cb)
	result3 := lo.RepeatBy(5, cb)

	is.Equal([]int{}, result1)
	is.Equal([]int{0, 1}, result2)
	is.Equal([]int{0, 1, 4, 9, 16}, result3)
}
