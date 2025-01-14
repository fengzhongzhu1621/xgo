package queue

import (
	"fmt"
	"testing"

	queue "github.com/duke-git/lancet/v2/datastructure/queue"
)

//func NewArrayQueue[T any](capacity int) *ArrayQueue[T]
//
//type ArrayQueue[T any] struct {
//	items    []T
//	head     int
//	tail     int
//	capacity int
//	size     int
//}

// func (q *ArrayQueue[T]) Data() []T

func TestArrayQueueData(t *testing.T) {
	q := queue.NewArrayQueue[int](5)
	fmt.Println(q.Data()) // []
}

// func (q *ArrayQueue[T]) Enqueue(item T) bool
func TestArrayQueueEnqueue(t *testing.T) {
	q := queue.NewArrayQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Data()) // 1,2,3
}

// func (q *ArrayQueue[T]) Dequeue() (T, bool)
func TestArrayQueueDequeue(t *testing.T) {
	q := queue.NewArrayQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Dequeue()) // 1
	fmt.Println(q.Data())    // 2,3
}

// func (q *ArrayQueue[T]) Front() T
func TestArrayQueueFront(t *testing.T) {
	q := queue.NewArrayQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Front()) // 1
	fmt.Println(q.Data())  // 1,2,3
}

// func (q *ArrayQueue[T]) Back() T
func TestArrayQueueBack(t *testing.T) {
	q := queue.NewArrayQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Back()) // 3
	fmt.Println(q.Data()) // 1,2,3
}

// func (q *ArrayQueue[T]) Size() int
func TestArrayQueueSize(t *testing.T) {
	q := queue.NewArrayQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Size()) // 3
}

// func (q *ArrayQueue[T]) IsEmpty() bool
func TestArrayQueueIsEmpty(t *testing.T) {
	q := queue.NewArrayQueue[int](5)
	fmt.Println(q.IsEmpty()) // true

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.IsEmpty()) // false
}

// func (q *ArrayQueue[T]) IsFull() bool
func TestArrayQueueIsFull(t *testing.T) {
	q := queue.NewArrayQueue[int](3)
	fmt.Println(q.IsFull()) // false

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.IsFull()) // true
}

// func (q *ArrayQueue[T]) Clear()
func TestArrayQueueClear(t *testing.T) {
	q := queue.NewArrayQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Clear()

	fmt.Println(q.IsEmpty()) // true
}

// func (q *ArrayQueue[T]) Contain(value T) bool
func TestArrayQueueContain(t *testing.T) {
	q := queue.NewArrayQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Contain(1)) // true
	fmt.Println(q.Contain(4)) // false
}
