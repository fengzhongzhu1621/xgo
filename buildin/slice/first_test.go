package slice

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1, ok1 := lo.First([]int{1, 2, 3})
	result2, ok2 := lo.First([]int{})

	is.Equal(result1, 1)
	is.Equal(ok1, true)
	is.Equal(result2, 0)
	is.Equal(ok2, false)
}

func TestFirstOrEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.FirstOrEmpty([]int{1, 2, 3})
	result2 := lo.FirstOrEmpty([]int{})
	result3 := lo.FirstOrEmpty([]string{})

	is.Equal(result1, 1)
	is.Equal(result2, 0)
	is.Equal(result3, "")
}

func TestFirstOr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.FirstOr([]int{1, 2, 3}, 63)
	result2 := lo.FirstOr([]int{}, 23)
	result3 := lo.FirstOr([]string{}, "test")

	is.Equal(result1, 1)
	is.Equal(result2, 23)
	is.Equal(result3, "test")
}
