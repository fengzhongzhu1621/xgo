package stack

import (
	"fmt"
	"testing"

	stack "github.com/duke-git/lancet/v2/datastructure/stack"
)

// LinkedStack is a stack implemented by linked list.

//type StackNode[T any] struct {
//	Value T
//	Next  *StackNode[T]
//}
//type LinkedStack[T any] struct {
//	top    *datastructure.StackNode[T]
//	length int
//}
//func NewLinkedStack[T any]() *LinkedStack[T]

// func (s *LinkedStack[T]) Push(value T)
func TestLinkedStackPush(t *testing.T) {
	sk := stack.NewLinkedStack[int]()
	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	fmt.Println(sk.Data()) //[]int{3, 2, 1}
}

// func (s *LinkedStack[T]) Pop() (*T, error)
func TestLinkedStackPop(t *testing.T) {
	sk := stack.NewLinkedStack[int]()
	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	val, err := sk.Pop()
	fmt.Println(err)  // nil
	fmt.Println(*val) // 3

	fmt.Println(sk.Data()) //[]int{2, 1}
}

// func (s *LinkedStack[T]) Peak() (*T, error)
func TestLinkedStackPeak(t *testing.T) {
	sk := stack.NewLinkedStack[int]()
	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	val, err := sk.Peak()
	fmt.Println(err)  // nil
	fmt.Println(*val) // 3

	fmt.Println(sk.Data()) //[]int{3, 2, 1}
}

// func (s *LinkedStack[T]) Data() []T
func TestLinkedStackData(t *testing.T) {
	sk := stack.NewLinkedStack[int]()
	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	fmt.Println(sk.Data()) //[]int{3, 2, 1}
}

// func (s *LinkedStack[T]) Size() int
func TestLinkedStackSize(t *testing.T) {
	sk := stack.NewLinkedStack[int]()
	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	fmt.Println(sk.Size()) // 3
}

// func (s *LinkedStack[T]) IsEmpty() bool
func TestLinkedStackIsEmpty(t *testing.T) {
	sk := stack.NewLinkedStack[int]()
	fmt.Println(sk.IsEmpty()) // true

	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	fmt.Println(sk.IsEmpty()) // false
}

// func (s *LinkedStack[T]) Clear()
func TestLinkedStackClear(t *testing.T) {
	sk := stack.NewLinkedStack[int]()

	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	sk.Clear()

	fmt.Println(sk.Data()) //[]int{}
}

// func (s *LinkedStack[T]) Print()
func TestLinkedStackPrint(t *testing.T) {
	sk := stack.NewLinkedStack[int]()

	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	sk.Print() //[ &{Value:3 Next:0xc000010260}, &{Value:2 Next:0xc000010250}, &{Value:1 Next:<nil>},  ]
}
