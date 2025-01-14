package hashmap

import (
	"fmt"
	"testing"

	heap "github.com/duke-git/lancet/v2/datastructure/hashmap"
)

// func (hm *HashMap) Iterate(iteratee func(key, value any))
func TestIterate(t *testing.T) {
	hm := heap.NewHashMap()
	hm.Put("a", 1)
	hm.Put("b", 2)
	hm.Put("c", 3)

	hm.Iterate(func(key, value any) {
		fmt.Println(key)
		fmt.Println(value)
	})
}
