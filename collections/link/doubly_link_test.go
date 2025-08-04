package link

import (
	"fmt"
	"testing"

	link "github.com/duke-git/lancet/v2/datastructure/link"
)

// DoublyLink is a linked list, whose node has a value, a next pointer points to next node and pre pointer points to previous node of the link.

//type LinkNode[T any] struct {
//	Value T
//	Pre   *LinkNode[T]
//	Next  *LinkNode[T]
//}
//type DoublyLink[T any] struct {
//	Head   *datastructure.LinkNode[T]
//	length int
//}
//
//func NewDoublyLink[T any]() *DoublyLink[T]

// func (link *DoublyLink[T]) Values() []T
func TestDoublyLinkValues(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	fmt.Println(lk.Values()) //[]int{1, 2, 3}
}

// func (link *DoublyLink[T]) InsertAt(index int, value T)
func TestDoublyLinkInsertAt(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAt(1, 1) // do nothing

	lk.InsertAt(0, 1)
	lk.InsertAt(1, 2)
	lk.InsertAt(2, 3)
	lk.InsertAt(2, 4)

	fmt.Println(lk.Values()) //[]int{1, 2, 4, 3}
}

// func (link *DoublyLink[T]) InsertAtHead(value T)
func TestDoublyLinkInsertAtHead(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtHead(1)
	lk.InsertAtHead(2)
	lk.InsertAtHead(3)

	fmt.Println(lk.Values()) //[]int{3, 2, 1}
}

// func (link *DoublyLink[T]) InsertAtTail(value T)
func TestDoublyLinkInsertAtTail(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	fmt.Println(lk.Values()) //[]int{1, 2, 3}
}

// func (link *DoublyLink[T]) DeleteAt(index int)
func TestDoublyLinkDeleteAt(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)
	lk.InsertAtTail(4)

	lk.DeleteAt(3)

	fmt.Println(lk.Values()) //[]int{1, 2, 3}
}

// func (link *DoublyLink[T]) DeleteAtHead()
func TestDoublyLinkDeleteAtHead(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)
	lk.InsertAtTail(4)

	lk.DeleteAtHead()

	fmt.Println(lk.Values()) //[]int{2, 3, 4}
}

// func (link *DoublyLink[T]) DeleteAtTail() error
func TestDoublyLinkDeleteAtTail(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	lk.DeleteAtTail()

	fmt.Println(lk.Values()) //[]int{1, 2}
}

// func (link *DoublyLink[T]) Reverse()
func TestDoublyLinkReverse(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	lk.Reverse()
	fmt.Println(lk.Values()) //[]int{3, 2, 1}
}

// func (link *DoublyLink[T]) GetMiddleNode() *datastructure.LinkNode[T]
func TestDoublyLinkGetMiddleNode(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	midNode := lk.GetMiddleNode()
	fmt.Println(midNode.Value) // 2
}

// func (link *DoublyLink[T]) Size() int
func TestDoublyLinkSize(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	fmt.Println(lk.Size()) // 3
}

// func (link *DoublyLink[T]) IsEmpty() bool
func TestDoublyLinkIsEmpty(t *testing.T) {
	lk := link.NewDoublyLink[int]()
	fmt.Println(lk.IsEmpty()) // true

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	fmt.Println(lk.IsEmpty()) // false
}

// func (link *DoublyLink[T]) Clear()
func TestDoublyLinkClear(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	lk.Clear()

	fmt.Println(lk.Values())
}

func TestDoublyLinkPrint(t *testing.T) {
	lk := link.NewDoublyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	lk.Print()
}
