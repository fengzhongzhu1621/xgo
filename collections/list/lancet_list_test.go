package list

import (
	"fmt"
	"testing"

	list "github.com/duke-git/lancet/v2/datastructure/list"
)

// List is a linear table, implemented with slice. NewList function return a list pointer
//type List[T any] struct {
//	data []T
//}
//func NewList[T any](data []T) *List[T]

// func (l *List[T]) Contain(value T) bool
func TestContain(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})

	fmt.Println(li.Contain(1)) //true
	fmt.Println(li.Contain(0)) //false

	data := li.Data()

	fmt.Println(data) //[]int{1, 2, 3}
}

// Return the value pointer at index in list
// func (l *List[T]) ValueOf(index int) (*T, bool)
func TestValueOf(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})
	v, ok := li.ValueOf(0)

	fmt.Println(*v) //1
	fmt.Println(ok) //true
}

// Returns the index of value in the list. if not found return -1
// func (l *List[T]) IndexOf(value T) int
func TestIndexOf(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})

	fmt.Println(li.IndexOf(1)) //0
	fmt.Println(li.IndexOf(0)) //-1
}

// Returns the index of the last occurrence of the value in this list if not found return -1
// func (l *List[T]) LastIndexOf(value T) int
func TestLastIndexOf(t *testing.T) {
	li := list.NewList([]int{1, 2, 3, 1})

	fmt.Println(li.LastIndexOf(1)) // 3
	fmt.Println(li.LastIndexOf(0)) //-1
}

// IndexOfFunc returns the first index satisfying f(v). if not found return -1
// func (l *List[T]) IndexOfFunc(f func(T) bool) int
func TestIndexOfFunc(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})

	fmt.Println(li.IndexOfFunc(func(a int) bool { return a == 1 })) //0
	fmt.Println(li.IndexOfFunc(func(a int) bool { return a == 0 })) //-1
}

// LastIndexOfFunc returns the index of the last occurrence of the value in this list satisfying f(data[i]). if not found return -1
// func (l *List[T]) LastIndexOfFunc(f func(T) bool) int
func TestLastIndexOfFunc(t *testing.T) {
	li := list.NewList([]int{1, 2, 3, 1})

	fmt.Println(li.LastIndexOfFunc(func(a int) bool { return a == 1 })) // 3
	fmt.Println(li.LastIndexOfFunc(func(a int) bool { return a == 0 })) //-1
}

// func (l *List[T]) Push(value T)
func TestPush(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})
	li.Push(4)

	fmt.Println(li.Data()) //[]int{1, 2, 3, 4}
}

// func (l *List[T]) PopFirst() (*T, bool)
func TestPopFirst(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})
	v, ok := li.PopFirst()

	fmt.Println(*v)        //1
	fmt.Println(ok)        //true
	fmt.Println(li.Data()) //2, 3
}

// func (l *List[T]) PopLast() (*T, bool)
func TestPopLast(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})
	v, ok := li.PopLast()

	fmt.Println(*v)        //3
	fmt.Println(ok)        //true
	fmt.Println(li.Data()) //1, 2
}

// Delete the value of list at index, if index is not between 0 and length of list data, do nothing
// func (l *List[T]) DeleteAt(index int)
func TestDeleteAt(t *testing.T) {
	li := list.NewList([]int{1, 2, 3, 4})

	li.DeleteAt(-1)
	fmt.Println(li.Data()) //1,2,3,4

	li.DeleteAt(4)
	fmt.Println(li.Data()) //1,2,3,4

	li.DeleteAt(0)
	fmt.Println(li.Data()) //2,3,4

	li.DeleteAt(2)
	fmt.Println(li.Data()) //2,3
}

// Insert value into list at index, if index is not between 0 and length of list data, do nothing
// func (l *List[T]) InsertAt(index int, value T)
func TestInsertAt(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})

	li.InsertAt(-1, 0)
	fmt.Println(li.Data()) //1,2,3

	li.InsertAt(4, 0)
	fmt.Println(li.Data()) //1,2,3

	li.InsertAt(3, 4)
	fmt.Println(li.Data()) //1,2,3,4

	li.InsertAt(2, 4)
	fmt.Println(li.Data()) //1,2,4,3,4
}

// Update value of list at index, index should between 0 and list size - 1
// func (l *List[T]) UpdateAt(index int, value T)
func TestUpdateAt(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})

	li.UpdateAt(-1, 0)
	fmt.Println(li.Data()) //1,2,3

	li.UpdateAt(2, 4)
	fmt.Println(li.Data()) //1,2,4

	li.UpdateAt(3, 5)
	fmt.Println(li.Data()) //1,2,4
}

// Compare a list to another list, use reflect.DeepEqual on every element
// func (l *List[T]) Equal(other *List[T]) bool
func TestEqual(t *testing.T) {
	li1 := list.NewList([]int{1, 2, 3, 4})
	li2 := list.NewList([]int{1, 2, 3, 4})
	li3 := list.NewList([]int{1, 2, 3})

	fmt.Println(li1.Equal(li2)) //true
	fmt.Println(li1.Equal(li3)) //false
}

// func (l *List[T]) IsEmpty() bool
func TestIsEmpty(t *testing.T) {
	li1 := list.NewList([]int{1, 2, 3})
	li2 := list.NewList([]int{})

	fmt.Println(li1.IsEmpty()) //false
	fmt.Println(li2.IsEmpty()) //true
}

// func (l *List[T]) Clear()
func TestClear(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})
	li.Clear()

	fmt.Println(li.Data()) // empty
}

// func (l *List[T]) Clone() *List[T]
func TestClone(t *testing.T) {
	li := list.NewList([]int{1, 2, 3})
	cloneList := li.Clone()

	fmt.Println(cloneList.Data()) // 1,2,3
}

// Merge two list, return new list, don't change original list
// func (l *List[T]) Merge(other *List[T]) *List[T]
func TestMerge(t *testing.T) {
	li1 := list.NewList([]int{1, 2, 3, 4})
	li2 := list.NewList([]int{4, 5, 6})
	li3 := li1.Merge(li2)

	fmt.Println(li3.Data()) //1, 2, 3, 4, 4, 5, 6
}

// func (l *List[T]) Size() int
func TestSize(t *testing.T) {
	li := list.NewList([]int{1, 2, 3, 4})

	fmt.Println(li.Size()) //4
}

// func (l *List[T]) Cap() int
func TestCap(t *testing.T) {
	data := make([]int, 0, 100)

	li := list.NewList(data)

	fmt.Println(li.Cap()) // 100
}

// func (l *List[T]) Swap(i, j int)
func TestSwap(t *testing.T) {
	li := list.NewList([]int{1, 2, 3, 4})
	li.Swap(0, 3)

	fmt.Println(li.Data()) //4, 2, 3, 1
}

// func (l *List[T]) Reverse()
func TestReverse(t *testing.T) {
	li := list.NewList([]int{1, 2, 3, 4})
	li.Reverse()

	fmt.Println(li.Data()) //4, 3, 2, 1
}

// func (l *List[T]) Unique()
func TestUnique(t *testing.T) {
	li := list.NewList([]int{1, 2, 2, 3, 4})
	li.Unique()

	fmt.Println(li.Data()) //1,2,3,4
}

// func (l *List[T]) Union(other *List[T]) *List[T]
func TestUnion(t *testing.T) {
	li1 := list.NewList([]int{1, 2, 3, 4})
	li2 := list.NewList([]int{4, 5, 6})
	li3 := li1.Union(li2)

	fmt.Println(li3.Data()) //1,2,3,4,5,6
}

// func (l *List[T]) Intersection(other *List[T]) *List[T]
func TestIntersection(t *testing.T) {
	li1 := list.NewList([]int{1, 2, 3, 4})
	li2 := list.NewList([]int{4, 5, 6})
	li3 := li1.Intersection(li2)

	fmt.Println(li3.Data()) //4
}

// func (l *List[T]) Difference(other *List[T]) *List[T]
func TestDifference(t *testing.T) {
	list1 := list.NewList([]int{1, 2, 3})
	list2 := list.NewList([]int{1, 2, 4})

	list3 := list1.Difference(list2)

	fmt.Println(list3.Data()) //3
}

// func (l *List[T]) SymmetricDifference(other *List[T]) *List[T]
func TestSymmetricDifference(t *testing.T) {
	list1 := list.NewList([]int{1, 2, 3})
	list2 := list.NewList([]int{1, 2, 4})

	list3 := list1.SymmetricDifference(list2)

	fmt.Println(list3.Data()) //3, 4
}

// TestRetainAll 仅保留此列表中包含在给定列表中的元素。
// Retains only the elements in this list that are contained in the given list.
// func (l *List[T]) RetainAll(list *List[T]) bool
func TestRetainAll(t *testing.T) {
	list1 := list.NewList([]int{1, 2, 3, 4})
	list2 := list.NewList([]int{1, 2, 3, 4})
	list3 := list.NewList([]int{1, 2, 3, 4})

	retain := list.NewList([]int{1, 2})
	retain1 := list.NewList([]int{2, 3})
	retain2 := list.NewList([]int{1, 2, 5})

	list1.RetainAll(retain)
	list2.RetainAll(retain1)
	list3.RetainAll(retain2)

	fmt.Println(list1.Data()) //1, 2
	fmt.Println(list2.Data()) //2, 3
	fmt.Println(list3.Data()) //1, 2
}

// Removes from this list all of its elements that are contained in the given list.
// func (l *List[T]) DeleteAll(list *List[T]) bool
func TestDeleteAll(t *testing.T) {
	list1 := list.NewList([]int{1, 2, 3, 4})
	list2 := list.NewList([]int{1, 2, 3, 4})
	list3 := list.NewList([]int{1, 2, 3, 4})

	del := list.NewList([]int{1})
	del1 := list.NewList([]int{2, 3})
	del2 := list.NewList([]int{1, 2, 5})

	list1.DeleteAll(del)
	list2.DeleteAll(del1)
	list3.DeleteAll(del2)

	fmt.Println(list1.Data()) //2,3,4
	fmt.Println(list2.Data()) //1,4
	fmt.Println(list3.Data()) //3,4
}

// func (l *List[T]) ForEach(consumer func(T))
func TestForEach(t *testing.T) {
	list1 := list.NewList([]int{1, 2, 3, 4})

	result := make([]int, 0)
	list1.ForEach(func(i int) {
		result = append(result, i)
	})

	fmt.Println(result) //1,2,3,4
}

// Returns an iterator over the elements in this list in proper sequence.
// func (l *List[T]) Iterator() iterator.Iterator[T]
func TestIterator(t *testing.T) {
	list1 := list.NewList([]int{1, 2, 3, 4})

	iterator := list1.Iterator()

	result := make([]int, 0)
	for iterator.HasNext() {
		item, _ := iterator.Next()
		result = append(result, item)
	}

	fmt.Println(result) //1,2,3,4
}

// Converts a list to a map based on iteratee function.
// func ListToMap[T any, K comparable, V any](list *List[T], iteratee func(T) (K, V)) map[K]V
func TestListToMap(t *testing.T) {
	list1 := list.NewList([]int{1, 2, 3, 4})

	result := list.ListToMap(list1, func(n int) (int, bool) {
		return n, n > 1
	})

	fmt.Println(result) //map[int]bool{1: false, 2: true, 3: true, 4: true}
}

// SubList returns a sub list of the original list between the specified fromIndex, inclusive, and toIndex, exclusive.
// func (l *List[T]) SubList(fromIndex, toIndex int) *List[T]
func TestSubList(t *testing.T) {
	l := list.NewList([]int{1, 2, 3, 4, 5, 6})

	fmt.Println(l.SubList(2, 5)) // []int{3, 4, 5}
}

// DeleteIf delete all satisfying f(data[i]), returns count of removed elements
// func (l *List[T]) DeleteIf(f func(T) bool) int
func TestDeleteIf(t *testing.T) {
	l := list.NewList([]int{1, 1, 1, 1, 2, 3, 1, 1, 4, 1, 1, 1, 1, 1, 1})

	fmt.Println(l.DeleteIf(func(a int) bool { return a == 1 })) // 12
	fmt.Println(l.Data())                                       // []int{2, 3, 4}
}
