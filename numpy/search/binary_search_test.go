package sort

import (
	"testing"

	"github.com/duke-git/lancet/v2/algorithm"
	"github.com/stretchr/testify/assert"
)

type intComparator struct{}

func (c *intComparator) Compare(v1 any, v2 any) int {
	val1, _ := v1.(int)
	val2, _ := v2.(int)

	// ascending order
	if val1 < val2 {
		return -1
	} else if val1 > val2 {
		return 1
	}
	return 0
}

// TestBinarySearch 二分搜索在已排序的切片中搜索目标，递归调用自身。如果找到目标，则返回目标的索引。否则函数返回-1。
// func BinarySearch[T any](sortedSlice []T, target T, lowIndex, highIndex int, comparator constraints.Comparator) int
func TestBinarySearch(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
	comparator := &intComparator{}

	result1 := algorithm.BinarySearch(numbers, 5, 0, len(numbers)-1, comparator)
	result2 := algorithm.BinarySearch(numbers, 9, 0, len(numbers)-1, comparator)

	assert.Equal(t, 4, result1)
	assert.Equal(t, -1, result2)
}
