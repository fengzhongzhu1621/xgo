package list

import (
	"fmt"
	"testing"

	list "github.com/duke-git/lancet/v2/datastructure/list"
)

//type CopyOnWriteList[T any] struct {
//	data []T
//	lock sync.
//}
//
//func NewCopyOnWriteList() *CopyOnWriteList

// func (l *CopyOnWriteList[T]) Size() int
func TestSafeListSize(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	fmt.Println(l.Size()) // 3
}

// func (c *CopyOnWriteList[T]) Get(index int) *T
func TestSafeListGet(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	fmt.Println(*l.Get(2)) // 3
}

// func (c *CopyOnWriteList[T]) Set(index int, e T) (oldValue *T, ok bool)
func TestSafeListSet(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	fmt.Println(l.Set(2, 4))
}

// Returns the index of the value in the list, or -1 if not found.
// func (c *CopyOnWriteList[T]) IndexOf(e T) int
func TestSafeListIndexOf(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	fmt.Println(l.IndexOf(1))
}

// func (c *CopyOnWriteList[T]) LastIndexOf(e T) int
func TestSafeListLastIndexOf(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3, 1})
	fmt.Println(l.LastIndexOf(1))
}

// func (l *CopyOnWriteList[T]) IndexOfFunc(f func(T) bool) int
func TestSafeListIndexOfFunc(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})

	fmt.Println(l.IndexOfFunc(func(a int) bool { return a == 1 })) // 0
	fmt.Println(l.IndexOfFunc(func(a int) bool { return a == 0 })) //-1
}

// func (l *CopyOnWriteList[T]) LastIndexOfFunc(f func(T) bool) int
func TestSafeListLastIndexOfFunc(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3, 1})

	fmt.Println(l.LastIndexOfFunc(func(a int) bool { return a == 1 })) // 3
	fmt.Println(l.LastIndexOfFunc(func(a int) bool { return a == 0 })) //-1
}

// func (c *CopyOnWriteList[T]) IsEmpty() bool
func TestSafeListIsEmpty(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{})
	fmt.Println(l.IsEmpty())
}

// func (c *CopyOnWriteList[T]) Contain(e T) bool
func TestSafeListContain(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	fmt.Println(l.Contain(1))
}

// func (c *CopyOnWriteList[T]) ValueOf(index int) []T
func TestSafeListValueOf(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	fmt.Println(l.ValueOf(2))
}

// func (c *CopyOnWriteList[T]) Add(e T) bool
func TestSafeListAdd(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	l.Add(4)
	fmt.Println(l.Size()) // 4
}

// func (c *CopyOnWriteList[T]) AddAll(e []T) bool
func TestSafeListAddAll(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	l.AddAll([]int{4, 5, 6})
	fmt.Println(l.Size()) // 6
}

// func (c *CopyOnWriteList[T]) AddByIndex(index int, e T) bool
func TestSafeListAddByIndex(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	l.AddByIndex(2, 6)
	fmt.Println(l.Size()) // 4
}

// func (c *CopyOnWriteList[T]) DeleteAt(index int) (oldValue *T, ok bool)
func TestSafeListDeleteAt(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	l.DeleteAt(2)
	fmt.Println(l.Size()) // 3
}

// func (c *CopyOnWriteList[T]) DeleteIf(func(T) bool) (oldValue *T, ok bool)
func TestSafeListDeleteIf(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	l.DeleteIf(func(i int) bool {
		return i == 2
	})
	fmt.Println(l.Size()) // 2
}

// func (c *CopyOnWriteList[T]) DeleteBy(e T) (*T bool)
func TestSafeListDeleteBy(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3})
	l.DeleteBy(2)
	fmt.Println(l.Size()) // 2
}

// func (c *CopyOnWriteList[T]) DeleteRange(start int, end int)
func TestSafeListDeleteRange(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	l.DeleteRange(2, 5)
	fmt.Println(l) // [1 2 6 7 8 9]
}

// func (c *CopyOnWriteList[T]) Equal(e []T) bool
func TestSafeListEqual(t *testing.T) {
	l := list.NewCopyOnWriteList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	fmt.Println(l.Equal(&[]int{1, 2, 3, 4, 5, 6, 7, 8, 9}))
}
