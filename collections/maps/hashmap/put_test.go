package hashmap

import (
	"fmt"

	heap "github.com/duke-git/lancet/v2/datastructure/hashmap"
)

func main() {
	hm := heap.NewHashMap()
	hm.Put("a", 1)

	val := hm.Get("a")
	fmt.Println(val) // 1
}
