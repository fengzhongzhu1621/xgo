package slice

import (
	"reflect"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestToSlice Creates a slice of give items.
// func ToSlice[T any](items ...T) []T
func TestToSlice(t *testing.T) {
	result := slice.ToSlice("a", "b", "c")
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

// TestToSlicePointer Returns a pointer to the slices of a variable parameter transformation
// func ToSlicePointer[T any](items ...T) []*T
func TestToSlicePointer(t *testing.T) {
	str1 := "a"
	str2 := "b"

	result := slice.ToSlicePointer(str1, str2)
	expect := []*string{&str1, &str2}

	isEqual := reflect.DeepEqual(result, expect)
	assert.Equal(t, isEqual, isEqual)
}

// TestWithout Creates a slice excluding all given items.
// func Without[T comparable](slice []T, items ...T) []T
func TestWithout(t *testing.T) {
	result := slice.Without([]int{1, 2, 3, 4}, 1, 2)
	assert.Equal(t, []int{3, 4}, result)
}
