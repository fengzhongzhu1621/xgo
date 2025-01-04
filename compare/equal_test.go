package compare

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/compare"
)

// TestEqual 检查两个值是否相等（同时检查类型和值）。
// func Equal(left, right any) bool
func TestEqual(t *testing.T) {
	result1 := compare.Equal(1, 1)
	result2 := compare.Equal("1", "1")
	result3 := compare.Equal([]int{1, 2, 3}, []int{1, 2, 3})
	result4 := compare.Equal(map[int]string{1: "a", 2: "b"}, map[int]string{1: "a", 2: "b"})

	result5 := compare.Equal(1, "1")
	result6 := compare.Equal(1, int64(1))
	result7 := compare.Equal([]int{1, 2}, []int{1, 2, 3})

	fmt.Println(result1) // true
	fmt.Println(result2) // true
	fmt.Println(result3) // true
	fmt.Println(result4) // true
	fmt.Println(result5) // false
	fmt.Println(result6) // false
	fmt.Println(result7) // false
}
