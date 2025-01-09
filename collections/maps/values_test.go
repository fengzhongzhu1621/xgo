package maps

import (
	"fmt"
	"sort"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
)

// TestValues Returns a slice of the map's values.
// func Values[K comparable, V any](m map[K]V) []V
func TestValues(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "a",
		3: "b",
		4: "c",
		5: "d",
	}

	values := maputil.Values(m)
	sort.Strings(values)

	fmt.Println(values)

	// Output:
	// [a a b c d]
}

// TestValuesBy Creates a slice whose element is the result of function mapper invoked by every map's value.
// func ValuesBy[K comparable, V any, T any](m map[K]V, mapper func(item V) T) []T
func TestValuesBy(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}
	values := maputil.ValuesBy(m, func(v string) string {
		switch v {
		case "a":
			return "a-1"
		case "b":
			return "b-2"
		case "c":
			return "c-3"
		default:
			return ""
		}
	})

	sort.Strings(values)

	fmt.Println(values)

	// Output:
	// [a-1 b-2 c-3]
}
