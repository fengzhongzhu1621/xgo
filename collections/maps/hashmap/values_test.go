package hashmap

import (
	"fmt"
	"testing"

	heap "github.com/duke-git/lancet/v2/datastructure/hashmap"
)

// func (hm *HashMap) Values() []any
func TestValues(t *testing.T) {
	hm := heap.NewHashMap()
	hm.Put("a", 1)
	hm.Put("b", 2)
	hm.Put("c", 3)

	values := hm.Values()
	fmt.Println(values) //[]interface{2, 1, 3}
}

// func (hm *HashMap) FilterByValue(perdicate func(value any) bool) *HashMap
func TestFilterByValue(t *testing.T) {
	hm := heap.NewHashMap()

	hm.Put("a", 1)
	hm.Put("b", 2)
	hm.Put("c", 3)
	hm.Put("d", 4)
	hm.Put("e", 5)
	hm.Put("f", 6)

	filteredHM := hm.FilterByValue(func(value any) bool {
		return value.(int) == 1 || value.(int) == 3
	})

	fmt.Println(filteredHM.Size()) // 2
}
