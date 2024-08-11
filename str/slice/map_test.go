package slice

import (
	"testing"

	"github.com/araujo88/lambda-go/pkg/core"
	"github.com/stretchr/testify/assert"
)

// Map 函数对片段中的每个元素应用给定的函数，返回一个包含转换后元素的新片段。
func TestMap(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	doubled := core.Map(numbers, func(x int) int { return x * 2 })
	assert.Equal(t, doubled, []int{2, 4, 6, 8, 10})
}
