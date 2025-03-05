package lancet

import (
	"fmt"
	"testing"

	queue "github.com/duke-git/lancet/v2/datastructure/queue"
)

// Circular queue implemented by slice.
// Return a CircularQueue pointer with specific capacity

//func NewCircularQueue[T any](capacity int) *CircularQueue[T]
//
//type CircularQueue[T any] struct {
//	data  []T
//	front int
//	rear  int
//	capacity  int
//}

// func (q *CircularQueue[T]) Data() []T
func TestCircularQueueData(t *testing.T) {
	q := queue.NewCircularQueue[int](5)
	fmt.Println(q.Data()) // []
}

// func (q *CircularQueue[T]) Enqueue(value T) error
func TestCircularQueueEnqueue(t *testing.T) {
	q := queue.NewCircularQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Data()) // 1,2,3
}

// func (q *CircularQueue[T]) Dequeue() (*T, bool)
func TestCircularQueueDequeue(t *testing.T) {
	q := queue.NewCircularQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	val, _ := q.Dequeue()
	fmt.Println(*val)     // 1
	fmt.Println(q.Data()) // 2,3
}

// func (q *CircularQueue[T]) Front() T
func TestCircularQueueFront(t *testing.T) {
	q := queue.NewCircularQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Front()) // 1
	fmt.Println(q.Data())  // 1,2,3
}

// func (q *CircularQueue[T]) Back() T
func TestCircularQueueBack(t *testing.T) {
	q := queue.NewCircularQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Back()) // 3
	fmt.Println(q.Data()) // 1,2,3
}

// func (q *CircularQueue[T]) Size() int
func TestCircularQueueSize(t *testing.T) {
	q := queue.NewCircularQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Size()) // 3
}

// func (q *CircularQueue[T]) IsEmpty() bool
func TestCircularQueueIsEmpty(t *testing.T) {
	q := queue.NewCircularQueue[int](5)
	fmt.Println(q.IsEmpty()) // true

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.IsEmpty()) // false
}

// func (q *CircularQueue[T]) IsFull() bool
func TestCircularQueueIsFull(t *testing.T) {
	q := queue.NewCircularQueue[int](3)
	fmt.Println(q.IsFull()) // false

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.IsFull()) // true
}

// func (q *CircularQueue[T]) Clear()
func TestCircularQueueClear(t *testing.T) {
	q := queue.NewCircularQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Clear()

	fmt.Println(q.IsEmpty()) // true
}

// func (q *CircularQueue[T]) Contain(value T) bool
func TestCircularQueueContain(t *testing.T) {
	q := queue.NewCircularQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Contain(1)) // true
	fmt.Println(q.Contain(4)) // false
}
