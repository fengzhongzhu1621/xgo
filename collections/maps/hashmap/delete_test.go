package hashmap

import (
	"fmt"
	"testing"

	heap "github.com/duke-git/lancet/v2/datastructure/hashmap"
)

// // func (hm *HashMap) Delete(key any)
func TestDelete(t *testing.T) {
	hm := heap.NewHashMap()
	hm.Put("a", 1)
	val := hm.Get("a")
	fmt.Println(val) // 1

	hm.Delete("a")
	val = hm.Get("a")
	fmt.Println(val) // nil
}
