package slice

import (
	"reflect"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
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

func TestWithout2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Without([]int{0, 2, 10}, 0, 1, 2, 3, 4, 5)
	result2 := lo.Without([]int{0, 7}, 0, 1, 2, 3, 4, 5)
	result3 := lo.Without([]int{}, 0, 1, 2, 3, 4, 5)
	result4 := lo.Without([]int{0, 1, 2}, 0, 1, 2)
	result5 := lo.Without([]int{})
	is.Equal(result1, []int{10})
	is.Equal(result2, []int{7})
	is.Equal(result3, []int{})
	is.Equal(result4, []int{})
	is.Equal(result5, []int{})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.Without(allStrings, "")
	is.IsType(nonempty, allStrings, "type preserved")
}
