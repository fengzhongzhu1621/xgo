package hashmap

import (
	"fmt"
	"testing"

	hashmap "github.com/duke-git/lancet/v2/datastructure/hashmap"
)

// Make a HashMap instance with default capacity is 1 << 10.
// func NewHashMap() *HashMap
func TestNewHashMap(t *testing.T) {
	hm := hashmap.NewHashMap()
	fmt.Println(hm)
}

// func NewHashMapWithCapacity(size, capacity uint64) *HashMap
func TestNewHashMapWithCapacity(t *testing.T) {
	hm := hashmap.NewHashMapWithCapacity(uint64(100), uint64(1000))
	fmt.Println(hm)
}
