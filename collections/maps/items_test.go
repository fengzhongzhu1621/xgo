package maps

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/convertor"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/stretchr/testify/assert"
)

// func MapToSlice[T any, K comparable, V any](aMap map[K]V, iteratee func(K, V) T) []T
func TestMapToSlice(t *testing.T) {
	aMap := map[string]int{"a": 1, "b": 2, "c": 3}
	result := convertor.MapToSlice(aMap, func(key string, value int) string {
		return key + ":" + strconv.Itoa(value)
	})

	fmt.Println(result) //[]string{"a:1", "b:2", "c:3"}
}

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

// TestMapKeys Transforms a map to other type map by manipulating it's keys.
// func MapKeys[K comparable, V any, T comparable](m map[K]V, iteratee func(key K, value V) T) map[T]V
func TestMapKeys(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}

	result := maputil.MapKeys(m, func(k int, _ string) string {
		return strconv.Itoa(k)
	})

	fmt.Println(result)

	// Output:
	// map[1:a 2:b 3:c]
}

// TestMapValues Transforms a map to other type map by manipulating it's values.
// func MapValues[K comparable, V any, T any](m map[K]V, iteratee func(key K, value V) T) map[K]T
func TestMapValues(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}

	result := maputil.MapValues(m, func(k int, v string) string {
		return v + strconv.Itoa(k)
	})

	fmt.Println(result)

	// Output:
	// map[1:a1 2:b2 3:c3]
}

// TestTransform Transform a map to another type map.
// func Transform[K1 comparable, V1 any, K2 comparable, V2 any](m map[K1]V1, iteratee func(key K1, value V1) (K2, V2)) map[K2]V2
func TestTransform(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	result := maputil.Transform(m, func(k string, v int) (string, string) {
		return k, strconv.Itoa(v)
	})

	fmt.Println(result)

	// Output:
	// map[a:1 b:2 c:3]
}

// TestEntry Transforms a map into array of key/value pairs.
//
//	type Entry[K comparable, V any] struct {
//	    Key   K
//	    Value V
//	}
//
// func Entries[K comparable, V any](m map[K]V) []Entry[K, V]
func TesEntries(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	result := maputil.Entries(m)

	sort.Slice(result, func(i, j int) bool {
		return result[i].Value < result[j].Value
	})

	fmt.Println(result)

	// Output:
	// [{a 1} {b 2} {c 3}]
}
