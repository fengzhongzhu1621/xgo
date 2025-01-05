package slice

import (
	"testing"

	"github.com/araujo88/lambda-go/pkg/utils"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

func TestConcat(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}

	concatenated := utils.Concat(slice1, slice2)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, concatenated)
}

func TestLancetConcat(t *testing.T) {
	result1 := slice.Concat([]int{1, 2}, []int{3, 4})
	result2 := slice.Concat([]string{"a", "b"}, []string{"c"}, []string{"d"})

	assert.Equal(t, []int{1, 2, 3, 4}, result1)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result2)
}

// TestAppendIfAbsent If slice doesn't contain the item, append it to the slice.
// func AppendIfAbsent[T comparable](slice []T, item T) []T
func TestAppendIfAbsent(t *testing.T) {
	result1 := slice.AppendIfAbsent([]string{"a", "b"}, "b")
	result2 := slice.AppendIfAbsent([]string{"a", "b"}, "c")

	assert.Equal(t, []string{"a", "b"}, result1)
	assert.Equal(t, []string{"a", "b", "c"}, result2)
}
