package slice

import (
	"fmt"
	"strings"
	"testing"

	"github.com/araujo88/lambda-go/pkg/utils"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestUnique 去重切片
func TestUnique(t *testing.T) {
	withDuplicates := []int{1, 2, 2, 3, 3, 3, 4}
	unique := utils.Unique(withDuplicates)
	fmt.Println(unique) // Output: [1 2 3 4]
}

// TestLancetUnique Remove duplicate elements in slice.
// func Unique[T comparable](slice []T) []T
func TestLancetUnique(t *testing.T) {
	result := slice.Unique([]string{"a", "a", "b"})
	assert.Equal(t, []string{"a", "b"}, result)
}

// TestLancetUniqueBy Removes duplicate elements from the input slice based
// on the values returned by the iteratee function.
// this function maintains the order of the elements.
// func UniqueBy[T any, U comparable](slice []T, iteratee func(item T) U) []T
func TestLancetUniqueBy(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	result := slice.UniqueBy(nums, func(val int) int {
		// 对余数进行去重
		return val % 3
	})

	assert.Equal(t, []int{1, 2, 3}, result)
}

// TestUniqueByComparator Removes duplicate elements from the input slice using the provided comparator function.
// The function maintains the order of the elements.
// func UniqueByComparator[T comparable](slice []T, comparator func(item T, other T) bool) []T
func TestUniqueByComparator(t *testing.T) {
	uniqueNums := slice.UniqueByComparator([]int{1, 2, 3, 1, 2, 4, 5, 6, 4}, func(item int, other int) bool {
		return item == other
	})

	caseInsensitiveStrings := slice.UniqueByComparator([]string{"apple", "banana", "Apple", "cherry", "Banana", "date"}, func(item string, other string) bool {
		return strings.ToLower(item) == strings.ToLower(other)
	})

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

	result, err := slice.UniqueByField(users, "ID")
	if err != nil {
	}
	assert.Equal(t, []User{
		{ID: 1, Name: "a"},
		{ID: 2, Name: "b"},
	}, result)
}
