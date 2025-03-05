package lancet

import (
	"fmt"
	"testing"

	queue "github.com/duke-git/lancet/v2/datastructure/queue"
)

//func NewLinkedQueue[T any]() *LinkedQueue[T]
//
//type LinkedQueue[T any] struct {
//	head   *datastructure.QueueNode[T]
//	tail   *datastructure.QueueNode[T]
//	length int
//}
//type QueueNode[T any] struct {
//	Value T
//	Next  *QueueNode[T]
//}

// func (q *LinkedQueue[T]) Data() []T
func TestLinkedQueueData(t *testing.T) {
	q := queue.NewLinkedQueue[int]()
	fmt.Println(q.Data()) // []
}

// func (q *LinkedQueue[T]) Enqueue(value T)
func TestLinkedQueueEnqueue(t *testing.T) {
	q := queue.NewLinkedQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Data()) // 1,2,3
}

// func (q *LinkedQueue[T]) Dequeue() (T, error)
func TestLinkedQueueDequeue(t *testing.T) {
	q := queue.NewLinkedQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Dequeue()) // 1
	fmt.Println(q.Data())    // 2,3
}

// func (q *LinkedQueue[T]) Front() (*T, error)
func TestLinkedQueueFront(t *testing.T) {
	q := queue.NewLinkedQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Front()) // 1
	fmt.Println(q.Data())  // 1,2,3
}

// func (q *LinkedQueue[T]) Back() (*T, error)
func TestLinkedQueueBack(t *testing.T) {
	q := queue.NewLinkedQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Back()) // 3
	fmt.Println(q.Data()) // 1,2,3
}

// func (q *LinkedQueue[T]) Size() int
func TestLinkedQueueSize(t *testing.T) {
	q := queue.NewLinkedQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Size()) // 3
}

// func (q *LinkedQueue[T]) IsEmpty() bool
func TestLinkedQueueIsEmpty(t *testing.T) {
	q := queue.NewLinkedQueue[int]()
	fmt.Println(q.IsEmpty()) // true

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.IsEmpty()) // false
}

// func (q *LinkedQueue[T]) Clear()
func TestLinkedQueueClear(t *testing.T) {
	q := queue.NewLinkedQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Clear()

	fmt.Println(q.IsEmpty()) // true
}

// func (q *LinkedQueue[T]) Contain(value T) bool
func TestLinkedQueueContain(t *testing.T) {
	q := queue.NewLinkedQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Contain(1)) // true
	fmt.Println(q.Contain(4)) // false
}
