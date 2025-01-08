package maps

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/stretchr/testify/assert"
)

// TestForEach Executes iteratee funcation for every key and value pair in map.
// func ForEach[K comparable, V any](m map[K]V, iteratee func(key K, value V))
func TestForEach(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	var sum int

	maputil.ForEach(m, func(_ string, value int) {
		sum += value
	})

	assert.Equal(t, sum, 10)
}

// TestFilter Iterates over map, return a new map contains all key and value pairs pass the predicate function.
// func Filter[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) map[K]V
func TestFilter(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}

	isEven := func(_ string, value int) bool {
		return value%2 == 0
	}

	result := maputil.Filter(m, isEven)

	fmt.Println(result)

	// Output:
	// map[b:2 d:4]
}

// TestFilterByKeys Iterates over map, return a new map whose keys are all given keys.
// func FilterByKeys[K comparable, V any](m map[K]V, keys []K) map[K]V
func TestFilterByKeys(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}

	result := maputil.FilterByKeys(m, []string{"a", "b"})

	fmt.Println(result)

	// Output:
	// map[b:2 d:4]
}

// TestFilterByValues Iterates over map, return a new map whose values are all given values.
// func FilterByValues[K comparable, V comparable](m map[K]V, values []V) map[K]V
func TestFilterByValues(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}

	result := maputil.FilterByValues(m, []int{3, 4})

	fmt.Println(result)

	// Output:
	// map[b:2 d:4]
}

// TestOmitBy is the opposite of Filter, removes all the map elements for which the predicate function returns true.
// func OmitBy[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) map[K]V
func TestOmitBy(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	isEven := func(_ string, value int) bool {
		return value%2 == 0
	}

	result := maputil.OmitBy(m, isEven)

	fmt.Println(result)

	// Output:
	// map[a:1 c:3 e:5]
}

// TestOmitByKeys The opposite of FilterByKeys, extracts all the map elements which keys are not omitted.
// func OmitByKeys[K comparable, V any](m map[K]V, keys []K) map[K]V
func TestOmitByKeys(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}

	result := maputil.OmitByKeys(m, []string{"a", "b"})

	fmt.Println(result)

	// Output:
	// map[c:3 d:4 e:5]
}

// TestOmitByValues The opposite of FilterByValues. remov all elements whose value are in the give slice.
// func OmitByValues[K comparable, V comparable](m map[K]V, values []V) map[K]V
func TestOmitByValues(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}

	result := maputil.OmitByValues(m, []int{4, 5})

	fmt.Println(result)

	// Output:
	// map[a:1 b:2 c:3]
}
