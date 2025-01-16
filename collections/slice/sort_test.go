package slice

import (
	"reflect"
	"testing"

	"github.com/araujo88/lambda-go/pkg/sortgroup"
	"github.com/duke-git/lancet/v2/slice"

	"github.com/stretchr/testify/assert"
)

// TestSort Sorts a slice of any ordered type(number or string),
// use quick sort algrithm. Default sort order is ascending (asc),
// if want descending order, set param `sortOrder` to `desc`. Ordered type: number(all ints uints floats) or string.
// func Sort[T constraints.Ordered](slice []T, sortOrder ...string)
func TestSort(t *testing.T) {
	{
		numbers := []int{1, 4, 3, 2, 5}
		slice.Sort(numbers)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, numbers)

		slice.Sort(numbers, "desc")
		assert.Equal(t, []int{5, 4, 3, 2, 1}, numbers)

		strings := []string{"a", "d", "c", "b", "e"}
		slice.Sort(strings)
		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, strings)

		slice.Sort(strings, "desc")
		assert.Equal(t, []string{"e", "d", "c", "b", "a"}, strings)
	}

	{
		tests := []struct {
			name  string
			slice []int
			want  []int
		}{
			{"sort positive numbers", []int{5, 3, 4, 1, 2}, []int{1, 2, 3, 4, 5}},
			{"sort including negatives", []int{-1, -3, 2, 1, 0}, []int{-3, -1, 0, 1, 2}},
			{"empty slice", []int{}, []int{}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sortgroup.SortSlice(tt.slice) // Modify the slice in-place
				if !reflect.DeepEqual(tt.slice, tt.want) {
					t.Errorf("sortSlice() resulted in %v, want %v", tt.slice, tt.want)
				}
			})
		}
	}
}

// TestSortBy Sorts the slice in ascending order as determined by the less function. This sort is not guaranteed to be stable.
// func SortBy[T any](slice []T, less func(a, b T) bool)
func TestSortBy(t *testing.T) {
	numbers := []int{1, 4, 3, 2, 5}

	slice.SortBy(numbers, func(a, b int) bool {
		return a < b
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, numbers)

	type User struct {
		Name string
		Age  uint
	}

	users := []User{
		{Name: "a", Age: 21},
		{Name: "b", Age: 15},
		{Name: "c", Age: 100}}

	slice.SortBy(users, func(a, b User) bool {
		return a.Age < b.Age
	})

	assert.Equal(t, []User{
		{Name: "b", Age: 15},
		{Name: "a", Age: 21},
		{Name: "c", Age: 100},
	}, users)
}

// TestSortByField Sort struct slice by field.
// Slice element should be struct, `field` param type should be int, uint, string, or bool.
// Default sort type is ascending (asc), if descending order, set `sortType` param to desc.
// func SortByField(slice any, field string, sortType ...string) error
func TestSortByField(t *testing.T) {
	type User struct {
		Name string
		Age  uint
	}

	users := []User{
		{Name: "a", Age: 21},
		{Name: "b", Age: 15},
		{Name: "c", Age: 100},
	}

	err := slice.SortByField(users, "Age", "desc")
	if err != nil {
		return
	}

	assert.Equal(t, []User{
		{Name: "c", Age: 100},
		{Name: "a", Age: 21},
		{Name: "b", Age: 15},
	}, users)
}
