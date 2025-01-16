package slice

import (
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
