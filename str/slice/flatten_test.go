package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestFlatten Flatten slice with one level.
// func Flatten(slice any) any
func TestFlatten(t *testing.T) {
	arrs := [][][]string{{{"a", "b"}}, {{"c", "d"}}}

	result1 := slice.Flatten(arrs)
	assert.Equal(t, [][]string{{"a", "b"}, {"c", "d"}}, result1)

	result2 := slice.Flatten(result1)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result2)
}

// TestFlattenDeep flattens slice recursive.
// func FlattenDeep(slice any) any
func TestFlattenDeep(t *testing.T) {
	arrs := [][][]string{{{"a", "b"}}, {{"c", "d"}}}

	result := slice.FlattenDeep(arrs)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result)
}
