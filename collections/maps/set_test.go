package maps

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
)

// TestIntersect Iterates over maps, return a new map of key and value pairs in all given maps.
// func Intersect[K comparable, V any](maps ...map[K]V) map[K]V
func TestIntersect(t *testing.T) {
	m1 := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	m2 := map[string]int{
		"a": 1,
		"b": 2,
		"c": 6,
		"d": 7,
	}

	m3 := map[string]int{
		"a": 1,
		"b": 9,
		"e": 9,
	}

	result1 := maputil.Intersect(m1)
	result2 := maputil.Intersect(m1, m2)
	result3 := maputil.Intersect(m1, m2, m3)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// map[a:1 b:2 c:3]
	// map[a:1 b:2]
	// map[a:1]
}

// TestMinus Creates an map of whose key in mapA but not in mapB.
// func Minus[K comparable, V any](mapA, mapB map[K]V) map[K]V
func TestMinus(t *testing.T) {
	m1 := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	m2 := map[string]int{
		"a": 11,
		"b": 22,
		"d": 33,
	}

	result := maputil.Minus(m1, m2)

	fmt.Println(result)

	// Output:
	// map[c:3]
}

// TestIsDisjoint Checks two maps are disjoint if they have no keys in common.
// func IsDisjoint[K comparable, V any](mapA, mapB map[K]V) bool
func TestIsDisjoint(t *testing.T) {
	m1 := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	m2 := map[string]int{
		"d": 22,
	}

	m3 := map[string]int{
		"a": 22,
	}

	result1 := maputil.IsDisjoint(m1, m2)
	result2 := maputil.IsDisjoint(m1, m3)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// true
	// false
}
