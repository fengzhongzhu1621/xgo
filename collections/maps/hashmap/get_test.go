package hashmap

import (
	"fmt"
	"testing"

	heap "github.com/duke-git/lancet/v2/datastructure/hashmap"
)

// func (hm *HashMap) Get(key any) any
func TestGet(t *testing.T) {
	hm := heap.NewHashMap()
	val := hm.Get("a")

	fmt.Println(val) //nil
}
