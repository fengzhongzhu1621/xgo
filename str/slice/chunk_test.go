package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestChunk 按照 size 参数均分 slice
// func Chunk[T any](slice []T, size int) [][]T
func TestChunk(t *testing.T) {
	arr := []string{"a", "b", "c", "d", "e"}

	result1 := slice.Chunk(arr, 1)
	result2 := slice.Chunk(arr, 2)
	result3 := slice.Chunk(arr, 3)
	result4 := slice.Chunk(arr, 4)
	result5 := slice.Chunk(arr, 5)
	result6 := slice.Chunk(arr, 6)

	assert.Equal(t, [][]string{{"a"}, {"b"}, {"c"}, {"d"}, {"e"}}, result1)
	assert.Equal(t, [][]string{{"a", "b"}, {"c", "d"}, {"e"}}, result2)
	assert.Equal(t, [][]string{{"a", "b", "c"}, {"d", "e"}}, result3)
	assert.Equal(t, [][]string{{"a", "b", "c", "d"}, {"e"}}, result4)
	assert.Equal(t, [][]string{{"a", "b", "c", "d", "e"}}, result5)
	assert.Equal(t, [][]string{{"a", "b", "c", "d", "e"}}, result6)
}
