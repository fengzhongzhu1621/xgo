package set

import (
	"fmt"
	"testing"

	set "github.com/duke-git/lancet/v2/datastructure/set"
)

// type Set[T comparable] map[T]struct{}
// func New[T comparable](items ...T) Set[T]

// func FromSlice[T comparable](items []T) Set[T]
func TestToSlice(t *testing.T) {
	st := set.FromSlice([]int{1, 2, 2, 3})
	fmt.Println(st.ToSlice()) //1,2,3
}

// func (s Set[T]) ToSortedSlice() (v T, ok bool)
func TestToSortedSlice(t *testing.T) {
	s1 := set.New(1, 2, 3, 4, 5)
	type Person struct {
		Name string
		Age  int
	}
	s2 := set.FromSlice([]Person{{"Tom", 20}, {"Jerry", 18}, {"Spike", 25}})

	res1 := s1.ToSortedSlice(func(v1, v2 int) bool {
		return v1 < v2
	})

	res2 := s2.ToSortedSlice(func(v1, v2 Person) bool {
		return v1.Age < v2.Age
	})

	fmt.Println(res1) // [1 2 3 4 5]
	fmt.Println(res2) // [{Jerry 18} {Tom 20} {Spike 25}]
}

// func (s Set[T]) Add(items ...T)
func TestAdd(t *testing.T) {
	st := set.New[int]()
	st.Add(1, 2, 3)

	fmt.Println(st.ToSlice()) //1,2,3
}

// func (s Set[T]) AddIfNotExist(item T) bool
func TestAddIfNotExist(t *testing.T) {
	st := set.New[int]()
	st.Add(1, 2, 3)

	r1 := st.AddIfNotExist(1)
	r2 := st.AddIfNotExist(4)

	fmt.Println(r1)           // false
	fmt.Println(r2)           // true
	fmt.Println(st.ToSlice()) // 1,2,3,4
}

// func (s Set[T]) AddIfNotExistBy(item T, checker func(element T) bool) bool
func TestAddIfNotExistBy(t *testing.T) {
	st := set.New[int]()
	st.Add(1, 2)

	ok := st.AddIfNotExistBy(3, func(val int) bool {
		return val%2 != 0
	})
	fmt.Println(ok) // true

	notOk := st.AddIfNotExistBy(4, func(val int) bool {
		return val%2 != 0
	})
	fmt.Println(notOk) // false

	fmt.Println(st.ToSlice()) // 1, 2, 3
}

// func (s Set[T]) Delete(items ...T)
func TestDelete(t *testing.T) {
	st := set.New[int]()
	st.Add(1, 2, 3)

	st.Delete(3)
	fmt.Println(st.ToSlice()) //1,2
}

// func (s Set[T]) Contain(item T) bool
func TestContain(t *testing.T) {
	st := set.New[int]()
	st.Add(1, 2, 3)

	fmt.Println(st.Contain(1)) //true
	fmt.Println(st.Contain(4)) //false
}

// func (s Set[T]) ContainAll(other Set[T]) bool
func TestContainAll(t *testing.T) {
	set1 := set.New(1, 2, 3)
	set2 := set.New(1, 2)
	set3 := set.New(1, 2, 3, 4)

	fmt.Println(set1.ContainAll(set2)) //true
	fmt.Println(set1.ContainAll(set3)) //false
}

// func (s Set[T]) Size() int
func TestSize(t *testing.T) {
	set1 := set.New(1, 2, 3)

	fmt.Println(set1.Size()) //3
}

// func (s Set[T]) Clone() Set[T]
func TestClone(t *testing.T) {
	set1 := set.New(1, 2, 3)
	set2 := set1.Clone()

	fmt.Println(set1.Size() == set2.Size()) //true
	fmt.Println(set1.ContainAll(set2))      //true
}

// func (s Set[T]) Equal(other Set[T]) bool
func TestEqual(t *testing.T) {
	set1 := set.New(1, 2, 3)
	set2 := set.New(1, 2, 3)
	set3 := set.New(1, 2, 3, 4)

	fmt.Println(set1.Equal(set2)) //true
	fmt.Println(set1.Equal(set3)) //false
}

// func (s Set[T]) Iterate(fn func(item T))
func TestIterate(t *testing.T) {
	set1 := set.New(1, 2, 3)
	arr := []int{}
	set1.Iterate(func(item int) {
		arr = append(arr, item)
	})

	fmt.Println(arr) //1,2,3
}

// func (s Set[T]) EachWithBreak(iteratee func(item T) bool)
func TestEachWithBreak(t *testing.T) {
	s := set.New(1, 2, 3, 4, 5)

	var sum int

	s.EachWithBreak(func(n int) bool {
		if n > 3 {
			return false
		}
		sum += n
		return true
	})

	fmt.Println(sum) //6
}

// func (s Set[T]) IsEmpty() bool
func TestIsEmpty(t *testing.T) {
	set1 := set.New(1, 2, 3)
	set2 := set.New[int]()

	fmt.Println(set1.IsEmpty()) //false
	fmt.Println(set2.IsEmpty()) //true
}

// func (s Set[T]) Union(other Set[T]) Set[T]
func TestUnion(t *testing.T) {
	set1 := set.New(1, 2, 3)
	set2 := set.New(2, 3, 4, 5)
	set3 := set1.Union(set2)

	fmt.Println(set3.ToSlice()) //1,2,3,4,5
}

// func (s Set[T]) Intersection(other Set[T]) Set[T]
func TestIntersection(t *testing.T) {
	set1 := set.New(1, 2, 3)
	set2 := set.New(2, 3, 4, 5)
	set3 := set1.Intersection(set2)

	fmt.Println(set3.ToSlice()) //2,3
}

// func (s Set[T]) SymmetricDifference(other Set[T]) Set[T]
func TestSymmetricDifference(t *testing.T) {
	set1 := set.New(1, 2, 3)
	set2 := set.New(2, 3, 4, 5)
	set3 := set1.SymmetricDifference(set2)

	fmt.Println(set3.ToSlice()) //1,4,5
}

// Create an set of whose element in origin set but not in compared set
// func (s Set[T]) Minus(comparedSet Set[T]) Set[T]
func TestMinus(t *testing.T) {
	set1 := set.New(1, 2, 3)
	set2 := set.New(2, 3, 4, 5)
	set3 := set.New(2, 3)

	res1 := set1.Minus(set2)
	fmt.Println(res1.ToSlice()) //1

	res2 := set2.Minus(set3)
	fmt.Println(res2.ToSlice()) //4,5
}

// Delete the top element of set then return it, if set is empty, return nil-value of T and false.
// func (s Set[T]) Pop() (v T, ok bool)
func TestPop(t *testing.T) {
	s := set.New[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	val, ok := s.Pop()

	fmt.Println(val) // 3
	fmt.Println(ok)  // true
}
