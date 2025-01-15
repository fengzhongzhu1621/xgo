package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestCompact 创建一些切片，其中删除了所有假值。值false、nil、0和""是假值。
// func Compact[T comparable](slice []T) []T
func TestCompact(t *testing.T) {
	result1 := slice.Compact([]int{0})
	result2 := slice.Compact([]int{0, 1, 2, 3})
	result3 := slice.Compact([]string{"", "a", "b", "0"})
	result4 := slice.Compact([]bool{false, true, true})

	assert.Equal(t, []int{}, result1)
	assert.Equal(t, []int{1, 2, 3}, result2)
	assert.Equal(t, []string{"a", "b", "0"}, result3)
	assert.Equal(t, []bool{true, true}, result4)
}

func TestWithoutEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Compact([]int{0, 1, 2})
	result2 := lo.Compact([]int{1, 2})
	result3 := lo.Compact([]int{})
	is.Equal(result1, []int{1, 2})
	is.Equal(result2, []int{1, 2})
	is.Equal(result3, []int{})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.Compact(allStrings)
	is.IsType(nonempty, allStrings, "type preserved")
}
