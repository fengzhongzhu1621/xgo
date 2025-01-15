package maps

import (
	"fmt"
	"sort"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
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

func TestValues2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.Values(map[string]int{"foo": 1, "bar": 2})
	sort.Ints(r1)
	is.Equal(r1, []int{1, 2})

	r2 := lo.Values(map[string]int{})
	is.Empty(r2)

	r3 := lo.Values(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
	sort.Ints(r3)
	is.Equal(r3, []int{1, 2, 3})

	r4 := lo.Values[string, int]()
	is.Equal(r4, []int{})

	r5 := lo.Values(map[string]int{"foo": 1, "bar": 2}, map[string]int{"foo": 1, "bar": 3})
	sort.Ints(r5)
	is.Equal(r5, []int{1, 1, 2, 3})
}

func TestUniqValues(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.UniqValues(map[string]int{"foo": 1, "bar": 2})
	sort.Ints(r1)
	is.Equal(r1, []int{1, 2})

	r2 := lo.UniqValues(map[string]int{})
	is.Empty(r2)

	r3 := lo.UniqValues(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
	sort.Ints(r3)
	is.Equal(r3, []int{1, 2, 3})

	r4 := lo.UniqValues[string, int]()
	is.Equal(r4, []int{})

	r5 := lo.UniqValues(map[string]int{"foo": 1, "bar": 2}, map[string]int{"foo": 1, "bar": 3})
	sort.Ints(r5)
	is.Equal(r5, []int{1, 2, 3})

	r6 := lo.UniqValues(map[string]int{"foo": 1, "bar": 1}, map[string]int{"foo": 1, "bar": 3})
	sort.Ints(r6)
	is.Equal(r6, []int{1, 3})

	// check order
	r7 := lo.UniqValues(map[string]int{"foo": 1}, map[string]int{"bar": 3})
	is.Equal(r7, []int{1, 3})
}
