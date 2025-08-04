package heap

import (
	"fmt"
	"testing"

	heap "github.com/duke-git/lancet/v2/datastructure/heap"
)

// MaxHeap is a binary heap tree implemented by slice,
// The key of the root node is both greater than or equal to the key value of the left subtree and greater than or equal to the key value of the right subtree.
//
//	type MaxHeap[T any] struct {
//		data       []T
//		comparator constraints.Comparator
//	}
// func NewMaxHeap[T any](comparator constraints.Comparator) *MaxHeap[T]

type intComparator struct{}

func (c *intComparator) Compare(v1, v2 any) int {
	val1, _ := v1.(int)
	val2, _ := v2.(int)

	if val1 < val2 {
		return -1
	} else if val1 > val2 {
		return 1
	}
	return 0
}

// func (h *MaxHeap[T]) Push(value T)
func TestMaxHeapPush(t *testing.T) {
	maxHeap := heap.NewMaxHeap[int](&intComparator{})
	values := []int{6, 5, 2, 4, 7, 10, 12, 1, 3, 8, 9, 11}

	for _, v := range values {
		maxHeap.Push(v)
	}

	fmt.Println(maxHeap.Data()) //[]int{12, 9, 11, 4, 8, 10, 7, 1, 3, 5, 6, 2}
}

// func (h *MaxHeap[T]) Pop() (T, bool)
func TestMaxHeapPop(t *testing.T) {
	maxHeap := heap.NewMaxHeap[int](&intComparator{})
	values := []int{6, 5, 2, 4, 7, 10, 12, 1, 3, 8, 9, 11}

	for _, v := range values {
		maxHeap.Push(v)
	}
	val, ok := maxHeap.Pop()

	fmt.Println(val) // 12
	fmt.Println(ok)  // true
}

// func (h *MaxHeap[T]) Peek() (T, bool)
func TestMaxHeapPeek(t *testing.T) {
	maxHeap := heap.NewMaxHeap[int](&intComparator{})
	values := []int{6, 5, 2, 4, 7, 10, 12, 1, 3, 8, 9, 11}

	for _, v := range values {
		maxHeap.Push(v)
	}
	val, _ := maxHeap.Peek()

	fmt.Println(val)            // 12
	fmt.Println(maxHeap.Size()) // 12
}

// func (h *MaxHeap[T]) Data() []T
func TestMaxHeapData(t *testing.T) {
	maxHeap := heap.NewMaxHeap[int](&intComparator{})
	values := []int{6, 5, 2, 4, 7, 10, 12, 1, 3, 8, 9, 11}

	for _, v := range values {
		maxHeap.Push(v)
	}

	fmt.Println(maxHeap.Data()) //[]int{12, 9, 11, 4, 8, 10, 7, 1, 3, 5, 6, 2}
}

// func (h *MaxHeap[T]) Size() int
func TestMaxHeapSize(t *testing.T) {
	maxHeap := heap.NewMaxHeap[int](&intComparator{})
	values := []int{6, 5, 2}

	for _, v := range values {
		maxHeap.Push(v)
	}

	fmt.Println(maxHeap.Size()) // 3
}

// func (h *MaxHeap[T]) PrintStructure()
func TestMaxHeapPrintStructure(t *testing.T) {
	maxHeap := heap.NewMaxHeap[int](&intComparator{})
	values := []int{6, 5, 2, 4, 7, 10, 12, 1, 3, 8, 9, 11}

	for _, v := range values {
		maxHeap.Push(v)
	}

	maxHeap.PrintStructure()
	//        12
	//    9       11
	//  4   8   10   7
	// 1 3 5 6 2
}
