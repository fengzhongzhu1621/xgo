package slice

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestTimes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Times(3, func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	})

	is.Equal(len(result1), 3)
	is.Equal(result1, []string{"0", "1", "2"})
}

// TestToSlice Creates a slice of give items.
// func ToSlice[T any](items ...T) []T
func TestToSlice(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		{
			result1 := lo.MapToSlice(map[int]int{1: 5, 2: 6, 3: 7, 4: 8}, func(k int, v int) string {
				return fmt.Sprintf("%d_%d", k, v)
			})
			result2 := lo.MapToSlice(map[int]int{1: 5, 2: 6, 3: 7, 4: 8}, func(k int, _ int) string {
				return strconv.FormatInt(int64(k), 10)
			})

			is.Equal(len(result1), 4)
			is.Equal(len(result2), 4)
			is.ElementsMatch(result1, []string{"1_5", "2_6", "3_7", "4_8"})
			is.ElementsMatch(result2, []string{"1", "2", "3", "4"})
		}
	}

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
