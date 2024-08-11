package slice

import (
	"testing"

	"github.com/araujo88/lambda-go/pkg/core"
	"github.com/stretchr/testify/assert"
)

// Foldl 和 Foldr 是两个功能强大的函数，通过对每个元素应用一个函数和一个累加器，可以将切片还原为一个单一值。
// 这两个函数的区别在于它们遍历切片的方向。
func TestFold(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	sum := core.Foldl(func(acc, x int) int { return acc + x }, 0, numbers)
	assert.Equal(t, sum, 15)
}
