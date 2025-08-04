package slice

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/araujo88/lambda-go/pkg/utils"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gookit/goutil/arrutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestUnique 去重切片
// func Unique[T comparable](slice []T) []T
func TestUnique(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.Uniq([]int{1, 2, 2, 1})

		is.Equal(len(result1), 2)
		is.Equal(result1, []int{1, 2})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Uniq(allStrings)
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
		withDuplicates := []int{1, 2, 2, 3, 3, 3, 4}
		unique := utils.Unique(withDuplicates)
		fmt.Println(unique) // Output: [1 2 3 4]

		tests := []struct {
			name  string
			slice []int
			want  []int
		}{
			{"unique elements", []int{1, 2, 2, 3, 4, 4, 4, 5}, []int{1, 2, 3, 4, 5}},
			{"all unique", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
			{"empty slice", []int{}, []int{}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := utils.Unique(tt.slice); !reflect.DeepEqual(got, tt.want) {
					t.Errorf(
						"Unique() = %v, want %v, failed in test case: %s",
						got,
						tt.want,
						tt.name,
					)
				}
			})
		}
	}

	{
		result := slice.Unique([]string{"a", "a", "b"})
		assert.Equal(t, []string{"a", "b"}, result)
	}

	{
		assert.Equal(t, []int{2}, arrutil.Unique[int]([]int{2}))
		assert.Equal(t, []int{2, 3, 4}, arrutil.Unique[int]([]int{2, 3, 2, 4}))
		assert.Equal(t, []uint{2, 3, 4}, arrutil.Unique([]uint{2, 3, 2, 4}))
		assert.Equal(
			t,
			[]string{"ab", "bc", "cd"},
			arrutil.Unique([]string{"ab", "bc", "ab", "cd"}),
		)
	}
}

// TestUniqueBy Removes duplicate elements from the input slice based
// on the values returned by the iteratee function.
// this function maintains the order of the elements.
// func UniqueBy[T any, U comparable](slice []T, iteratee func(item T) U) []T
func TestUniqueBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.UniqBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
			return i % 3
		})

		is.Equal(len(result1), 3)
		is.Equal(result1, []int{0, 1, 2})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.UniqBy(allStrings, func(i string) string {
			return i
		})
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
		nums := []int{1, 2, 3, 4, 5, 6}
		result := slice.UniqueBy(nums, func(val int) int {
			// 对余数进行去重
			return val % 3
		})

		assert.Equal(t, []int{1, 2, 3}, result)
	}
}

// TestUniqueByComparator Removes duplicate elements from the input slice using the provided comparator function.
// The function maintains the order of the elements.
// func UniqueByComparator[T comparable](slice []T, comparator func(item T, other T) bool) []T
func TestUniqueByComparator(t *testing.T) {
	uniqueNums := slice.UniqueByComparator(
		[]int{1, 2, 3, 1, 2, 4, 5, 6, 4},
		func(item int, other int) bool {
			return item == other
		},
	)

	caseInsensitiveStrings := slice.UniqueByComparator(
		[]string{"apple", "banana", "Apple", "cherry", "Banana", "date"},
		func(item string, other string) bool {
			// return strings.ToLower(item) == strings.ToLower(other)
			return strings.EqualFold(item, other)
		},
	)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, uniqueNums)
	assert.Equal(t, []string{"apple", "banana", "cherry", "date"}, caseInsensitiveStrings)
}

// TestUniqueByConcurrent Removes duplicate elements from the slice by parallel.
// func UniqueByConcurrent[T comparable](slice []T, comparator func(item T, other T) bool, numThreads int) []T
func TestUniqueByConcurrent(t *testing.T) {
	nums := []int{1, 2, 3, 1, 2, 4, 5, 6, 4, 7}
	comparator := func(item int, other int) bool { return item == other }

	result := slice.UniqueByConcurrent(nums, comparator, 4)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, result)
}

// TestUniqueByField Remove duplicate elements in struct slice by struct field.
// func UniqueByField[T any](slice []T, field string) ([]T, error)
func TestUniqueByField(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	users := []User{
		{ID: 1, Name: "a"},
		{ID: 2, Name: "b"},
		{ID: 1, Name: "c"},
	}

	result, _ := slice.UniqueByField(users, "ID")
	assert.Equal(t, []User{
		{ID: 1, Name: "a"},
		{ID: 2, Name: "b"},
	}, result)
}

func TestFindUniques(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.FindUniques([]int{1, 2, 3})

	is.Equal(3, len(result1))
	is.Equal([]int{1, 2, 3}, result1)

	result2 := lo.FindUniques([]int{1, 2, 2, 3, 1, 2})

	is.Equal(1, len(result2))
	is.Equal([]int{3}, result2)

	result3 := lo.FindUniques([]int{1, 2, 2, 1})

	is.Equal(0, len(result3))
	is.Equal([]int{}, result3)

	result4 := lo.FindUniques([]int{})

	is.Equal(0, len(result4))
	is.Equal([]int{}, result4)

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.FindUniques(allStrings)
	is.IsType(nonempty, allStrings, "type preserved")
}

func TestFindUniquesBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.FindUniquesBy([]int{0, 1, 2}, func(i int) int {
		return i % 3
	})

	is.Equal(3, len(result1))
	is.Equal([]int{0, 1, 2}, result1)

	result2 := lo.FindUniquesBy([]int{0, 1, 2, 3, 4}, func(i int) int {
		return i % 3
	})

	is.Equal(1, len(result2))
	is.Equal([]int{2}, result2)

	result3 := lo.FindUniquesBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
		return i % 3
	})

	is.Equal(0, len(result3))
	is.Equal([]int{}, result3)

	result4 := lo.FindUniquesBy([]int{}, func(i int) int {
		return i % 3
	})

	is.Equal(0, len(result4))
	is.Equal([]int{}, result4)

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.FindUniquesBy(allStrings, func(i string) string {
		return i
	})
	is.IsType(nonempty, allStrings, "type preserved")
}

func TestFindDuplicates(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.FindDuplicates([]int{1, 2, 2, 1, 2, 3})

	is.Equal(2, len(result1))
	is.Equal([]int{1, 2}, result1)

	result2 := lo.FindDuplicates([]int{1, 2, 3})

	is.Equal(0, len(result2))
	is.Equal([]int{}, result2)

	result3 := lo.FindDuplicates([]int{})

	is.Equal(0, len(result3))
	is.Equal([]int{}, result3)

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.FindDuplicates(allStrings)
	is.IsType(nonempty, allStrings, "type preserved")
}

func TestFindDuplicatesBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.FindDuplicatesBy([]int{3, 4, 5, 6, 7}, func(i int) int {
		return i % 3
	})

	is.Equal(2, len(result1))
	is.Equal([]int{3, 4}, result1)

	result2 := lo.FindDuplicatesBy([]int{0, 1, 2, 3, 4}, func(i int) int {
		return i % 5
	})

	is.Equal(0, len(result2))
	is.Equal([]int{}, result2)

	result3 := lo.FindDuplicatesBy([]int{}, func(i int) int {
		return i % 3
	})

	is.Equal(0, len(result3))
	is.Equal([]int{}, result3)

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.FindDuplicatesBy(allStrings, func(i string) string {
		return i
	})
	is.IsType(nonempty, allStrings, "type preserved")
}
