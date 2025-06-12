package buildin

// Go 1.18+ 引入的泛型比较功能，使用 cmp 包进行值比较

import (
	"cmp"
	"fmt"
	"testing"
)

func Min[T cmp.Ordered](a, b T) T {
	if cmp.Less(a, b) {
		return a
	}
	return b
}

func TestCmp(t *testing.T) {
	result := cmp.Compare(5, 10) // -1
	fmt.Println(result)
	result = cmp.Compare(10, 5) // 1
	fmt.Println(result)
	result = cmp.Compare(5, 5) // 0
	fmt.Println(result)

	isLess := cmp.Less(5, 10) // true
	fmt.Println(isLess)
	isLess = cmp.Less(10, 5) // false
	fmt.Println(isLess)
	isLess = cmp.Less(10, 10) // false
	fmt.Println(isLess)

	minInt := Min(3, 5)        // minInt is 3
	minString := Min("a", "b") // minString is "a"

	fmt.Println(minInt, minString) // 3 a
}
