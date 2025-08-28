package maps

import (
	"fmt"
	"maps"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/gookit/goutil/arrutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

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
	t.Parallel()
	is := assert.New(t)

	{
		r1 := lo.OmitBy(
			map[string]int{"foo": 1, "bar": 2, "baz": 3},
			func(key string, value int) bool {
				return value%2 == 1
			},
		)

		is.Equal(r1, map[string]int{"bar": 2})

		type myMap map[string]int
		before := myMap{"": 0, "foobar": 6, "baz": 3}
		after := lo.PickBy(before, func(key string, value int) bool { return true })
		is.IsType(after, before, "type preserved")
	}

	{
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

		fmt.Println(result) // map[a:1 c:3 e:5]
	}

	{
		m1 := map[int]string{1: "a", 2: "b", 3: "c", 4: "d"}
		maps.DeleteFunc(m1, func(k int, v string) bool {
			return k%2 == 0
		})
		fmt.Println(m1) // map[1:a 3:c]
	}
}

// TestOmitByKeys The opposite of FilterByKeys, extracts all the map elements which keys are not omitted.
// func OmitByKeys[K comparable, V any](m map[K]V, keys []K) map[K]V
func TestOmitByKeys(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		r1 := lo.OmitByKeys(
			map[string]int{"foo": 1, "bar": 2, "baz": 3},
			[]string{"foo", "baz", "qux"},
		)

		is.Equal(r1, map[string]int{"bar": 2})

		type myMap map[string]int
		before := myMap{"": 0, "foobar": 6, "baz": 3}
		after := lo.OmitByKeys(before, []string{"foobar", "baz"})
		is.IsType(after, before, "type preserved")
	}

	{
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
		// map[c:3 d:4 e:5]}
	}
}

// TestOmitByValues The opposite of FilterByValues. remov all elements whose value are in the give slice.
// func OmitByValues[K comparable, V comparable](m map[K]V, values []V) map[K]V
func TestOmitByValues(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{

		r1 := lo.OmitByValues(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})

		is.Equal(r1, map[string]int{"bar": 2})

		type myMap map[string]int
		before := myMap{"": 0, "foobar": 6, "baz": 3}
		after := lo.OmitByValues(before, []int{0, 3})
		is.IsType(after, before, "type preserved")
	}

	{
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
}

func TestPickBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.PickBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
		return value%2 == 1
	})

	is.Equal(r1, map[string]int{"foo": 1, "baz": 3})

	type myMap map[string]int
	before := myMap{"": 0, "foobar": 6, "baz": 3}
	after := lo.PickBy(before, func(key string, value int) bool { return true })
	is.IsType(after, before, "type preserved")
}

func TestPickByKeys(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.PickByKeys(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz", "qux"})

	is.Equal(r1, map[string]int{"foo": 1, "baz": 3})

	type myMap map[string]int
	before := myMap{"": 0, "foobar": 6, "baz": 3}
	after := lo.PickByKeys(before, []string{"foobar", "baz"})
	is.IsType(after, before, "type preserved")
}

func TestPickByValues(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.PickByValues(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})

	is.Equal(r1, map[string]int{"foo": 1, "baz": 3})

	type myMap map[string]int
	before := myMap{"": 0, "foobar": 6, "baz": 3}
	after := lo.PickByValues(before, []int{0, 3})
	is.IsType(after, before, "type preserved")
}

func TestMapToColumnSlice(t *testing.T) {
	list1 := []map[string]any{
		{"name": "tom", "age": 23},
		{"name": "john", "age": 34},
	}

	flatArr := arrutil.Column(list1, func(obj map[string]any) (val any, find bool) {
		return obj["age"], true
	})
	fmt.Println(flatArr) // [23 34]
	assert.NotEmpty(t, flatArr)
	assert.Contains(t, flatArr, 23)
}
