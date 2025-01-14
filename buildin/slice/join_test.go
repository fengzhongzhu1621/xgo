package slice

import (
	"strings"
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

// TestLancetConcatBy Concats the elements of a slice into a single value using the provided separator and connector function.
// func ConcatBy[T any](slice []T, sep T, connector func(T, T) T) T
func TestLancetConcatBy(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}

	sep := Person{Name: " | ", Age: 0}

	personConnector := func(a, b Person) Person {
		// "" + Alice
		// "" + Alice + " | " + Bob
		// "" + Alice + " | " + Bob + " | " + Charlie
		//
		// 0 + 30
		// 0 + 30 + 0 + 25
		// 0 + 30 + 0 + 25 + 0 + 35
		return Person{Name: a.Name + b.Name, Age: a.Age + b.Age}
	}

	result := slice.ConcatBy(people, sep, personConnector)

	assert.Equal(t, "Alice | Bob | Charlie", result.Name)
	assert.Equal(t, 90, result.Age)
}

// TestJoin Join the slice item with specify separator.
// func Join[T any](s []T, separator string) string
func TestJoin(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	result1 := slice.Join(nums, ",")
	result2 := slice.Join(nums, "-")

	assert.Equal(t, "1,2,3,4,5", result1)
	assert.Equal(t, "1-2-3-4-5", result2)
}

// TestJoinFunc Joins the slice elements into a single string with the given separator.
// func JoinFunc[T any](slice []T, sep string, transform func(T) T) string
func TestJoinFunc(t *testing.T) {
	result := slice.JoinFunc([]string{"a", "b", "c"}, ", ", func(s string) string {
		return strings.ToUpper(s)
	})

	assert.Equal(t, "A, B, C", result)
}

// TestAppendIfAbsent If slice doesn't contain the item, append it to the slice.
// func AppendIfAbsent[T comparable](slice []T, item T) []T
func TestAppendIfAbsent(t *testing.T) {
	result1 := slice.AppendIfAbsent([]string{"a", "b"}, "b")
	result2 := slice.AppendIfAbsent([]string{"a", "b"}, "c")

	assert.Equal(t, []string{"a", "b"}, result1)
	assert.Equal(t, []string{"a", "b", "c"}, result2)
}
