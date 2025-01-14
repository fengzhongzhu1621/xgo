package hashmap

import (
	"fmt"
	"testing"

	heap "github.com/duke-git/lancet/v2/datastructure/hashmap"
)

func TestContains(t *testing.T) {
	hm := heap.NewHashMap()
	hm.Put("a", 1)

	fmt.Println(hm.Contains("a")) //true
	fmt.Println(hm.Contains("b")) //false
}
