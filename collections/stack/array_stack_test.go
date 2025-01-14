package stack

import (
	"fmt"
	"testing"

	stack "github.com/duke-git/lancet/v2/datastructure/stack"
)

// ArrayStack is a stack implemented by slice.
//type ArrayStack[T any] struct {
//	data   []T
//	length int
//}
//func NewArrayStack[T any]() *ArrayStack[T]

// func (s *ArrayStack[T]) Push(value T)

func TestPush(t *testing.T) {
	sk := stack.NewArrayStack[int]()
	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	fmt.Println(sk.Data()) //[]int{3, 2, 1}
}

// func (s *ArrayStack[T]) Pop() (*T, error)
func TestPop(t *testing.T) {
	sk := stack.NewArrayStack[int]()
	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	val, err := sk.Pop()
	fmt.Println(err)  //nil
	fmt.Println(*val) //3

	fmt.Println(sk.Data()) //[]int{2, 1}
}

// func (s *ArrayStack[T]) Peak() (*T, error)
func TestPeak(t *testing.T) {
	sk := stack.NewArrayStack[int]()
	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	val, err := sk.Peak()
	fmt.Println(err)  //nil
	fmt.Println(*val) //3

	fmt.Println(sk.Data()) //[]int{3, 2, 1}
}

// func (s *ArrayStack[T]) Data() []T
func TestData(t *testing.T) {
	sk := stack.NewArrayStack[int]()
	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	fmt.Println(sk.Data()) //[]int{3, 2, 1}
}

// func (s *ArrayStack[T]) Size() int
func TestSize(t *testing.T) {
	sk := stack.NewArrayStack[int]()
	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	fmt.Println(sk.Size()) //3
}

// func (s *ArrayStack[T]) IsEmpty() bool
func TestIsEmpty(t *testing.T) {
	sk := stack.NewArrayStack[int]()
	fmt.Println(sk.IsEmpty()) //true

	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	fmt.Println(sk.IsEmpty()) //false
}

// func (s *ArrayStack[T]) Clear()
func TestClear(t *testing.T) {
	sk := stack.NewArrayStack[int]()

	sk.Push(1)
	sk.Push(2)
	sk.Push(3)

	sk.Clear()

	fmt.Println(sk.Data()) //[]int{}
}
