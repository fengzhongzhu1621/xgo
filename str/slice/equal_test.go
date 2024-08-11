package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// 检查两个切片是否相等，相等条件：切片长度相同，元素顺序和值都相同。
func TestEqual(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{1, 3, 2}

	result1 := slice.Equal(s1, s2)
	result2 := slice.Equal(s1, s3)

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}
