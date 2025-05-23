package slice

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gookit/goutil/arrutil"

	"github.com/araujo88/lambda-go/pkg/utils"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

func TestConcat(t *testing.T) {
	{
		slice1 := []int{1, 2, 3}
		slice2 := []int{4, 5, 6}

		concatenated := utils.Concat(slice1, slice2)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, concatenated)

		tests := []struct {
			name   string
			slice1 []int
			slice2 []int
			want   []int
		}{
			{"concat two slices", []int{1, 2, 3}, []int{4, 5}, []int{1, 2, 3, 4, 5}},
			{"concat empty with non-empty", []int{}, []int{1, 2, 3}, []int{1, 2, 3}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := utils.Concat(tt.slice1, tt.slice2); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Concat() = %v, want %v", got, tt.want)
				}
			})
		}
	}

	{
		result1 := slice.Concat([]int{1, 2}, []int{3, 4})
		result2 := slice.Concat([]string{"a", "b"}, []string{"c"}, []string{"d"})

		assert.Equal(t, []int{1, 2, 3, 4}, result1)
		assert.Equal(t, []string{"a", "b", "c", "d"}, result2)
	}
}

// TestConcatBy Concats the elements of a slice into a single value using the provided separator and connector function.
// func ConcatBy[T any](slice []T, sep T, connector func(T, T) T) T
func TestConcatBy(t *testing.T) {
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

func TestStringsToString(t *testing.T) {
	is := assert.New(t)

	is.Equal("a,b", arrutil.JoinStrings(",", []string{"a", "b"}...))
	is.Equal("a,b", arrutil.StringsJoin(",", []string{"a", "b"}...))
	is.Equal("a,b", arrutil.StringsJoin(",", "a", "b"))
}

func TestJoinTyped(t *testing.T) {
	assert.Equal(t, "", arrutil.JoinTyped[any](","))
	assert.Equal(t, "", arrutil.JoinTyped[any](",", nil))
	assert.Equal(t, "1,2", arrutil.JoinTyped(",", 1, 2))
	assert.Equal(t, "a,b", arrutil.JoinTyped(",", "a", "b"))
	assert.Equal(t, "1,a", arrutil.JoinTyped[any](",", 1, "a"))
}

func TestJoinSlice(t *testing.T) {
	assert.Equal(t, "", arrutil.JoinSlice(","))
	assert.Equal(t, "", arrutil.JoinSlice(",", nil))
	assert.Equal(t, "a,23,b", arrutil.JoinSlice(",", "a", 23, "b"))
}
