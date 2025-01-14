package hashmap

import (
	"fmt"
	"testing"

	heap "github.com/duke-git/lancet/v2/datastructure/hashmap"
)

// func (hm *HashMap) Keys() []any
func TestKeys(t *testing.T) {
	hm := heap.NewHashMap()
	hm.Put("a", 1)
	hm.Put("b", 2)
	hm.Put("c", 3)

	keys := hm.Keys()
	fmt.Println(keys) //[]interface{"a", "b", "c"}
}
