package lancet

import (
	"fmt"
	"testing"

	queue "github.com/duke-git/lancet/v2/datastructure/queue"
)

// Common queue implemented by slice.
// Return a PriorityQueue pointer with specific capacity, param `comparator` is used to compare values of type T in the queue.

//func NewPriorityQueue[T any](capacity int, comparator constraints.Comparator) *PriorityQueue[T]
//
//type PriorityQueue[T any] struct {
//	items      []T
//	size       int
//	comparator constraints.Comparator
//}

// func (q *PriorityQueue[T]) Enqueue(item T) bool
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

// Put element into queue, if queue is full, return false
// func (q *PriorityQueue[T]) Data() []T
func TestPriorityQueueData(t *testing.T) {
	comparator := &intComparator{}
	q := queue.NewPriorityQueue[int](10, comparator)
	fmt.Println(q.Data()) // []
}

// func (q *PriorityQueue[T]) Enqueue(item T) bool
func TestPriorityQueueEnqueue(t *testing.T) {
	comparator := &intComparator{}
	q := queue.NewPriorityQueue[int](10, comparator)
	for i := 1; i < 11; i++ {
		q.Enqueue(i)
	}

	fmt.Println(q.Data()) // 10, 9, 6, 7, 8, 2, 5, 1, 4, 3
}

// Remove head element of queue and return it
// func (q *PriorityQueue[T]) Dequeue() (T, bool)
func TestPriorityQueueDequeue(t *testing.T) {
	comparator := &intComparator{}
	q := queue.NewPriorityQueue[int](10, comparator)
	for i := 1; i < 11; i++ {
		q.Enqueue(i)
	}
	val, _ := q.Dequeue()
	fmt.Println(val)      // 10
	fmt.Println(q.Data()) // 9, 8, 6, 7, 3, 2, 5, 1, 4
}

// func (q *PriorityQueue[T]) IsEmpty() bool
func TestPriorityQueueIsEmpty(t *testing.T) {
	comparator := &intComparator{}
	q := queue.NewPriorityQueue[int](10, comparator)
	fmt.Println(q.IsEmpty()) // true

	for i := 1; i < 11; i++ {
		q.Enqueue(i)
	}
	fmt.Println(q.IsEmpty()) // false
}

// func (q *PriorityQueue[T]) IsFull() bool
func TestPriorityQueueIsFull(t *testing.T) {
	comparator := &intComparator{}
	q := queue.NewPriorityQueue[int](10, comparator)
	fmt.Println(q.IsFull()) // false

	for i := 1; i < 11; i++ {
		q.Enqueue(i)
	}
	fmt.Println(q.IsFull()) // true
}

// func (q *PriorityQueue[T]) Size() int
func TestPriorityQueueSize(t *testing.T) {
	comparator := &intComparator{}
	q := queue.NewPriorityQueue[int](10, comparator)
	fmt.Println(q.IsFull()) // false

	for i := 1; i < 5; i++ {
		q.Enqueue(i)
	}
	fmt.Println(q.Size()) // 4
}
