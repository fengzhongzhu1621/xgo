package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
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
