package sort

import (
	"testing"

	"github.com/duke-git/lancet/v2/algorithm"
	"github.com/stretchr/testify/assert"
)

// TestSelectionSort 使用选择排序算法对切片进行排序。参数comparator应实现constraints.Comparator。
// func SelectionSort[T any](slice []T, comparator constraints.Comparator)
func TestSelectionSort(t *testing.T) {
	numbers := []int{2, 1, 5, 3, 6, 4}
	comparator := &intComparator{}

	algorithm.SelectionSort(numbers, comparator)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, numbers)
}
