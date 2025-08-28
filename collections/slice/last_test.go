package slice

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestLast(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1, ok1 := lo.Last([]int{1, 2, 3})
	result2, ok2 := lo.Last([]int{})

	is.Equal(result1, 3)
	is.True(ok1)
	is.Equal(result2, 0)
	is.False(ok2)
}

func TestLastOrEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.LastOrEmpty([]int{1, 2, 3})
	result2 := lo.LastOrEmpty([]int{})
	result3 := lo.LastOrEmpty([]string{})

	is.Equal(result1, 3)
	is.Equal(result2, 0)
	is.Equal(result3, "")
}

func TestLastOr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.LastOr([]int{1, 2, 3}, 63)
	result2 := lo.LastOr([]int{}, 23)
	result3 := lo.LastOr([]string{}, "test")

	is.Equal(result1, 3)
	is.Equal(result2, 23)
	is.Equal(result3, "test")
}
