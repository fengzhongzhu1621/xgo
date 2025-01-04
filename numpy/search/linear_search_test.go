package sort

import (
	"testing"

	"github.com/duke-git/lancet/v2/algorithm"
	"github.com/stretchr/testify/assert"
)

// TestLinearSearch 根据相等函数返回目标在切片中的索引。如果找到目标，则返回目标的索引。否则，函数返回-1。
// func LinearSearch[T any](slice []T, target T, equal func(a, b T) bool) int
func TestLinearSearch(t *testing.T) {
	numbers := []int{3, 4, 5, 3, 2, 1}
	equalFunc := func(a, b int) bool {
		return a == b
	}

	result1 := algorithm.LinearSearch(numbers, 3, equalFunc)
	result2 := algorithm.LinearSearch(numbers, 6, equalFunc)

	assert.Equal(t, 0, result1)
	assert.Equal(t, -1, result2)
}
