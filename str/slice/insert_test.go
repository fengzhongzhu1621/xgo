package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestInsertAt insert the element into slice at index.
// func InsertAt[T any](slice []T, index int, value any) []T
func TestInsertAt(t *testing.T) {
	result1 := slice.InsertAt([]string{"a", "b", "c"}, 0, "1")
	result2 := slice.InsertAt([]string{"a", "b", "c"}, 1, "1")
	result3 := slice.InsertAt([]string{"a", "b", "c"}, 2, "1")
	result4 := slice.InsertAt([]string{"a", "b", "c"}, 3, "1")
	result5 := slice.InsertAt([]string{"a", "b", "c"}, 0, []string{"1", "2", "3"})

	assert.Equal(t, []string{"1", "a", "b", "c"}, result1)
	assert.Equal(t, []string{"a", "1", "b", "c"}, result2)
	assert.Equal(t, []string{"a", "b", "1", "c"}, result3)
	assert.Equal(t, []string{"a", "b", "c", "1"}, result4)
	assert.Equal(t, []string{"1", "2", "3", "a", "b", "c"}, result5)
}
