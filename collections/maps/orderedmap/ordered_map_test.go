package maps

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
)

// TestNewOrderedMap  Map Creates a new OrderedMap.
// func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V]
// func (om *OrderedMap[K, V]) Set(key K, value V)
// func (om *OrderedMap[K, V]) Get(key K) (V, bool)
// func (om *OrderedMap[K, V]) Delete(key K)
// func (om *OrderedMap[K, V]) Clear()
// func (om *OrderedMap[K, V]) Keys() []K
// func (om *OrderedMap[K, V]) Values() []V
// func (om *OrderedMap[K, V]) Elements() []struct
func TestNewOrderedMap(t *testing.T) {
	om := maputil.NewOrderedMap[string, int]()

	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	elements := om.Elements()
	fmt.Println(elements)
	// Output:
	// [{a 1} {b 2} {c 3}]

	val1, ok := om.Get("a")
	fmt.Println(val1, ok)

	val2, ok := om.Get("d")
	fmt.Println(val2, ok)

	// Output:
	// 1 true
	// 0 false

	om.Delete("b")
	fmt.Println(om.Keys())
	// Output:
	// [a c]

	om.Clear()
	fmt.Println(om.Keys())
	// Output:
	// []
}

// TestOrderedMapFront  Returns the first key-value pair.
//
//	func (om *OrderedMap[K, V]) Front() (struct {
//	    Key   K
//	    Value V
//	}, bool)
func TestOrderedMapFront(t *testing.T) {
	om := maputil.NewOrderedMap[string, int]()

	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	frontElement, ok := om.Front()
	fmt.Println(frontElement)
	fmt.Println(ok)

	// Output:
	// {a 1}
	// true
}

// TestOrderedMapFront  Returns the last key-value pair.
//
//	func (om *OrderedMap[K, V]) Back() (struct {
//	    Key   K
//	    Value V
//	}, bool)
func TestOrderedMapBack(t *testing.T) {
	om := maputil.NewOrderedMap[string, int]()

	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	backElement, ok := om.Back()
	fmt.Println(backElement)
	fmt.Println(ok)

	// Output:
	// {c 3}
	// true
}

// TestOrderedMapRange  Calls the given function for each key-value pair.
//
// func (om *OrderedMap[K, V]) Range(iteratee func(key K, value V) bool)
func TestOrderedMapRange(t *testing.T) {
	om := maputil.NewOrderedMap[string, int]()

	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	om.Range(func(key string, value int) bool {
		fmt.Println(key, value)
		return true
	})

	// Output:
	// a 1
	// b 2
	// c 3
}

// TestOrderedMapLen Returns the number of key-value pairs.
//
// func (om *OrderedMap[K, V]) Len() int
func TestOrderedMapLen(t *testing.T) {
	om := maputil.NewOrderedMap[string, int]()

	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	om.Len()

	fmt.Println(om.Len())

	// Output:
	// 3
}

// TestOrderedMapContains Returns true if the given key exists.
//
// func (om *OrderedMap[K, V]) Contains(key K) bool
func TestOrderedMapContains(t *testing.T) {
	om := maputil.NewOrderedMap[string, int]()

	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	result1 := om.Contains("a")
	result2 := om.Contains("d")

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// true
	// false
}

// TestOrderedMapIter Returns a channel that yields key-value pairs in order.
//
//	func (om *OrderedMap[K, V]) Iter() <-chan struct {
//	    Key   K
//	    Value V
//	}
func TestOrderedMapIter(t *testing.T) {
	om := maputil.NewOrderedMap[string, int]()

	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	for elem := range om.Iter() {
		fmt.Println(elem)
	}

	// Output:
	// {a 1}
	// {b 2}
	// {c 3}
}

// TestOrderedMapReverseIter Returns a channel that yields key-value pairs in reverse order.
//
//	func (om *OrderedMap[K, V]) ReverseIter() <-chan struct {
//	    Key   K
//	    Value V
//	}
func TestOrderedMapReverseIter(t *testing.T) {
	om := maputil.NewOrderedMap[string, int]()

	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	for elem := range om.ReverseIter() {
		fmt.Println(elem)
	}

	// Output:
	// {c 3}
	// {b 2}
	// {a 1}
}

// TestOrderedMapSortByKey Sorts the map by key given less function.
//
// func (om *OrderedMap[K, V]) SortByKey(less func(a, b K) bool)
func TestOrderedMapSortByKey(t *testing.T) {
	om := maputil.NewOrderedMap[int, string]()

	om.Set(3, "c")
	om.Set(1, "a")
	om.Set(4, "d")
	om.Set(2, "b")

	fmt.Println(om.Elements())
	// Output:
	// [{3 c} {1 a} {4 d} {2 b}]

	om.SortByKey(func(a, b int) bool {
		return a < b
	})

	fmt.Println(om.Elements())
	// Output:
	// [{1 a} {2 b} {3 c} {4 d}]
}

// TestOrderedMarshalJSON Implements the json.Marshaler interface.
// func (om *OrderedMap[K, V]) MarshalJSON() ([]byte, error)
func TestOrderedMarshalJSON(t *testing.T) {
	om := maputil.NewOrderedMap[int, string]()

	om.Set(3, "c")
	om.Set(1, "a")
	om.Set(4, "d")
	om.Set(2, "b")

	b, _ := om.MarshalJSON()

	fmt.Println(string(b))

	// Output:
	// {"a":1,"b":2,"c":3}
}

// TestOrderedUnmarshalJSON Implements the json.Unmarshaler interface.
// func (om *OrderedMap[K, V]) UnmarshalJSON(data []byte) error
func TestOrderedUnmarshalJSON(t *testing.T) {
	om := maputil.NewOrderedMap[string, int]()

	data := []byte(`{"a":1,"b":2,"c":3}`)

	om.UnmarshalJSON(data)

	fmt.Println(om.Elements())

	// Output:
	// [{a 1} {b 2} {c 3}]
}
