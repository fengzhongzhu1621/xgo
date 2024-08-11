package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// 判断 slice 是否包含指定值
func TestContain(t *testing.T) {
	result1 := slice.Contain([]string{"a", "b", "c"}, "a")
	result2 := slice.Contain([]int{1, 2, 3}, 4)

	assert.Equal(t, result1, true)
	assert.Equal(t, result2, false)
}
