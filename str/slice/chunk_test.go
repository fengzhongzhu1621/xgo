package slice

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
)

// 按照 size 参数均分 slice
func TestChunk(t *testing.T) {
	arr := []string{"a", "b", "c", "d", "e"}

	result1 := slice.Chunk(arr, 1)
	result2 := slice.Chunk(arr, 2)

	fmt.Println(result1)
	fmt.Println(result2)
}
