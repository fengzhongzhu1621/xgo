package slice

import (
	"github.com/duke-git/lancet/v2/slice"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsAscending Checks if a slice is ascending order.
// func IsAscending[T constraints.Ordered](slice []T) bool
func TestIsAscending(t *testing.T) {
	result1 := slice.IsAscending([]int{1, 2, 3, 4, 5})
	result2 := slice.IsAscending([]int{5, 4, 3, 2, 1})
	result3 := slice.IsAscending([]int{2, 1, 3, 4, 5})

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, false, result3)
}

// TestIsDescending Checks if a slice is descending order.
// func IsDescending[T constraints.Ordered](slice []T) bool
func TestIsDescending(t *testing.T) {
	result1 := slice.IsDescending([]int{5, 4, 3, 2, 1})
	result2 := slice.IsDescending([]int{1, 2, 3, 4, 5})
	result3 := slice.IsDescending([]int{2, 1, 3, 4, 5})

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
	assert.Equal(t, false, result3)
}

// TestIsSorted Checks if a slice is sorted (ascending or descending).
// func IsSorted[T constraints.Ordered](slice []T) bool
func TestIsSorted(t *testing.T) {
	result1 := slice.IsSorted([]int{5, 4, 3, 2, 1})
	result2 := slice.IsSorted([]int{1, 2, 3, 4, 5})
	result3 := slice.IsSorted([]int{2, 1, 3, 4, 5})

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, false, result3)
}

// TestIsSortedByKey Checks if a slice is sorted by iteratee function.
// func IsSortedByKey[T any, K constraints.Ordered](slice []T, iteratee func(item T) K) bool
func TestIsSortedByKey(t *testing.T) {
	result1 := slice.IsSortedByKey([]string{"a", "ab", "abc"}, func(s string) int {
		return len(s)
	})
	result2 := slice.IsSortedByKey([]string{"abc", "ab", "a"}, func(s string) int {
		return len(s)
	})
	result3 := slice.IsSortedByKey([]string{"abc", "a", "ab"}, func(s string) int {
		return len(s)
	})

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, false, result3)
}

// TestSort Sorts a slice of any ordered type(number or string),
// use quick sort algrithm. Default sort order is ascending (asc),
// if want descending order, set param `sortOrder` to `desc`. Ordered type: number(all ints uints floats) or string.
// func Sort[T constraints.Ordered](slice []T, sortOrder ...string)
func TestSort(t *testing.T) {
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
