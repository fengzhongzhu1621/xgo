package maps

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
)

// TestSortByKey Sorts the map by its keys and returns a new map with sorted keys.
//
// func SortByKey[K constraints.Ordered, V any](m map[K]V) (sortedKeysMap map[K]V)
func TestSortByKey(t *testing.T) {
	m := map[int]string{
		3: "c",
		1: "a",
		4: "d",
		2: "b",
	}

	result := maputil.SortByKey(m, func(a, b int) bool {
		return a < b
	})

	fmt.Println(result)

	// Output:
	// map[1:a 2:b 3:c 4:d]
}

// TestGetOrDefault returns the value of the given key or a default value if the key is not present.
//
// func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V
func TestGetOrDefault(t *testing.T) {
	m := map[int]string{
		3: "c",
		1: "a",
		4: "d",
		2: "b",
	}

	result1 := maputil.GetOrDefault(m, 1, "default")
	result2 := maputil.GetOrDefault(m, 6, "default")

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// a
	// default
}

// TestToSortedSlicesDefault Translate the key and value of the map into two slices that are sorted in ascending order according to the key’s value,
// with the position of the elements in the value slice corresponding to the key.
//
// func ToSortedSlicesDefault[K constraints.Ordered, V any](m map[K]V) ([]K, []V)
func TestToSortedSlicesDefault(t *testing.T) {
	m := map[int]string{
		1: "a",
		3: "c",
		2: "b",
	}

	keys, values := maputil.ToSortedSlicesDefault(m)

	fmt.Println(keys)
	fmt.Println(values)

	// Output:
	// [1 2 3]
	// [a b c]
}

// TestToSortedSlicesWithComparator Translate the key and value of the map into two slices that are sorted according to a custom sorting rule defined by a comparator function based on the key's value,
// with the position of the elements in the value slice corresponding to the key.
//
// func ToSortedSlicesWithComparator[K comparable, V any](m map[K]V, comparator func(a, b K) bool) ([]K, []V)
func TestToSortedSlicesWithComparator(t *testing.T) {
	m1 := map[time.Time]string{
		time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC): "today",
		time.Date(2024, 3, 30, 0, 0, 0, 0, time.UTC): "yesterday",
		time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC):  "tomorrow",
	}

	keys1, values1 := maputil.ToSortedSlicesWithComparator(m1, func(a, b time.Time) bool {
		return a.Before(b)
	})

	m2 := map[int]string{
		1: "a",
		3: "c",
		2: "b",
	}
	keys2, values2 := maputil.ToSortedSlicesWithComparator(m2, func(a, b int) bool {
		return a > b
	})

	fmt.Println(keys2)
	fmt.Println(values2)

	fmt.Println(keys1)
	fmt.Println(values1)

	// Output:
	// [3 2 1]
	// [c b a]
	// [2024-03-30 00:00:00 +0000 UTC 2024-03-31 00:00:00 +0000 UTC 2024-04-01 00:00:00 +0000 UTC]
	// [yesterday today tomorrow]
}
