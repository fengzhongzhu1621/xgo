package slice

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestFlatten Flatten slice with one level.
// func Flatten(slice any) any
func TestFlatten(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.Flatten([][]int{{0, 1}, {2, 3, 4, 5}})

		is.Equal(result1, []int{0, 1, 2, 3, 4, 5})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Flatten([]myStrings{allStrings})
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
		// flattens slice with one level.
		arrs := [][][]string{{{"a", "b"}}, {{"c", "d"}}}

		result1 := slice.Flatten(arrs)
		assert.Equal(t, [][]string{{"a", "b"}, {"c", "d"}}, result1)

		result2 := slice.Flatten(result1)
		assert.Equal(t, []string{"a", "b", "c", "d"}, result2)
	}
}

// TestFlattenDeep flattens slice recursive.
// func FlattenDeep(slice any) any
func TestFlattenDeep(t *testing.T) {
	is := assert.New(t)

	{
		tests := []struct {
			name        string
			collections [][]int
			want        []int
		}{
			{
				"nil",
				[][]int{nil},
				[]int{},
			},
			{
				"empty",
				[][]int{},
				[]int{},
			},
			{
				"empties",
				[][]int{{}, {}},
				[]int{},
			},
			{
				"same length",
				[][]int{{1, 3, 5}, {2, 4, 6}},
				[]int{1, 2, 3, 4, 5, 6},
			},
			{
				"different length",
				[][]int{{1, 3, 5, 6}, {2, 4}},
				[]int{1, 2, 3, 4, 5, 6},
			},
			{
				"many slices",
				[][]int{{1}, {2, 5, 8}, {3, 6}, {4, 7, 9, 10}},
				[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := lo.Interleave(tt.collections...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Interleave() = %v, want %v", got, tt.want)
				}
			})
		}

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Interleave(allStrings)
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
		arrs := [][][]string{{{"a", "b"}}, {{"c", "d"}}}

		result := slice.FlattenDeep(arrs)
		assert.Equal(t, []string{"a", "b", "c", "d"}, result)
	}
}

// TestFlatMap Manipulates a slice and transforms and flattens it to a slice of another type.
// func FlatMap[T any, U any](slice []T, iteratee func(index int, item T) []U) []U
func TestFlatMap(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.FlatMap([]int{0, 1, 2, 3, 4}, func(x int, _ int) []string {
			return []string{"Hello"}
		})
		result2 := lo.FlatMap([]int64{0, 1, 2, 3, 4}, func(x int64, _ int) []string {
			result := make([]string, 0, x)
			for i := int64(0); i < x; i++ {
				result = append(result, strconv.FormatInt(x, 10))
			}
			return result
		})

		is.Equal(len(result1), 5)
		is.Equal(len(result2), 10)
		is.Equal(result1, []string{"Hello", "Hello", "Hello", "Hello", "Hello"})
		is.Equal(result2, []string{"1", "2", "2", "3", "3", "3", "4", "4", "4", "4"})
	}

	{
		nums := []int{1, 2, 3, 4}

		result := slice.FlatMap(nums, func(i int, num int) []string {
			s := "hi-" + strconv.FormatInt(int64(num), 10)
			return []string{s}
		})

		assert.Equal(t, []string{"hi-1", "hi-2", "hi-3", "hi-4"}, result)
	}
}
