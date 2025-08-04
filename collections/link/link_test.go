package link

import (
	"fmt"
	"testing"

	link "github.com/duke-git/lancet/v2/datastructure/link"
)

//type LinkNode[T any] struct {
//	Value T
//	Next  *LinkNode[T]
//}
//type SinglyLink[T any] struct {
//	Head   *datastructure.LinkNode[T]
//	length int
//}
//func NewSinglyLink[T any]() *SinglyLink[T]

// func (link *SinglyLink[T]) Values() []T
func TestValues(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	fmt.Println(lk.Values()) //[]int{1, 2, 3}
}

// func (link *SinglyLink[T]) InsertAt(index int, value T)
func TestInsertAt(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAt(1, 1) // do nothing

	lk.InsertAt(0, 1)
	lk.InsertAt(1, 2)
	lk.InsertAt(2, 3)
	lk.InsertAt(2, 4)

	fmt.Println(lk.Values()) //[]int{1, 2, 4, 3}
}

// func (link *SinglyLink[T]) InsertAtHead(value T)
func TestInsertAtHead(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtHead(1)
	lk.InsertAtHead(2)
	lk.InsertAtHead(3)

	fmt.Println(lk.Values()) //[]int{3, 2, 1}
}

// func (link *SinglyLink[T]) InsertAtTail(value T)
func TestInsertAtTail(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	fmt.Println(lk.Values()) //[]int{1, 2, 3}
}

// func (link *SinglyLink[T]) DeleteAt(index int)
func TestDeleteAt(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)
	lk.InsertAtTail(4)

	lk.DeleteAt(3)

	fmt.Println(lk.Values()) //[]int{1, 2, 3}
}

// func (link *SinglyLink[T]) DeleteAtHead()
func TestDeleteAtHead(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)
	lk.InsertAtTail(4)

	lk.DeleteAtHead()

	fmt.Println(lk.Values()) //[]int{2, 3, 4}
}

// func (link *SinglyLink[T]) DeleteAtTail()
func TestDeleteAtTail(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	lk.DeleteAtTail()

	fmt.Println(lk.Values()) //[]int{1, 2}
}

// func (link *SinglyLink[T]) DeleteValue(value T)
func TestDeleteValue(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	lk.DeleteValue(2)
	fmt.Println(lk.Values()) //[]int{1, 3}
}

// func (link *SinglyLink[T]) Reverse()
func TestReverse(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	lk.Reverse()
	fmt.Println(lk.Values()) //[]int{3, 2, 1}
}

// func (link *SinglyLink[T]) GetMiddleNode() *datastructure.LinkNode[T]
func TestGetMiddleNode(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	midNode := lk.GetMiddleNode()
	fmt.Println(midNode.Value) // 2
}

// func (link *SinglyLink[T]) Size() int
func TestSize(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	fmt.Println(lk.Size()) // 3
}

// func (link *SinglyLink[T]) IsEmpty() bool
func TestIsEmpty(t *testing.T) {
	lk := link.NewSinglyLink[int]()
	fmt.Println(lk.IsEmpty()) // true

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	fmt.Println(lk.IsEmpty()) // false
}

// func (link *SinglyLink[T]) Clear()
func TestClear(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	lk.Clear()

	fmt.Println(lk.Values()) //
}

func TestPrint(t *testing.T) {
	lk := link.NewSinglyLink[int]()

	lk.InsertAtTail(1)
	lk.InsertAtTail(2)
	lk.InsertAtTail(3)

	lk.Print() //[ &{Value:1 Pre:<nil> Next:0xc0000a4048}, &{Value:2 Pre:<nil> Next:0xc0000a4060}, &{Value:3 Pre:<nil> Next:<nil>} ]
}
