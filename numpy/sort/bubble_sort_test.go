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

// TestBubbleSort 使用冒泡排序算法对切片进行排序。参数comparator应实现constraints.Comparator。
// func BubbleSort[T any](slice []T, comparator constraints.Comparator)
func TestBubbleSort(t *testing.T) {
	numbers := []int{2, 1, 5, 3, 6, 4}
	comparator := &intComparator{}

	algorithm.BubbleSort(numbers, comparator)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, numbers)
}
