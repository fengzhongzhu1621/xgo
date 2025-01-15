package maps

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/samber/lo"
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
	t.Parallel()
	is := assert.New(t)

	{
		r1 := lo.OmitBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
			return value%2 == 1
		})

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

		fmt.Println(result)

		// Output:
		// map[a:1 c:3 e:5]
	}

}

// TestOmitByKeys The opposite of FilterByKeys, extracts all the map elements which keys are not omitted.
// func OmitByKeys[K comparable, V any](m map[K]V, keys []K) map[K]V
func TestOmitByKeys(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		r1 := lo.OmitByKeys(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz", "qux"})

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

// TestMapKeys Transforms a map to other type map by manipulating it's keys.
// func MapKeys[K comparable, V any, T comparable](m map[K]V, iteratee func(key K, value V) T) map[T]V
func TestMapKeys(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.MapKeys(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
			return "Hello"
		})
		result2 := lo.MapKeys(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(_ int, v int) string {
			return strconv.FormatInt(int64(v), 10)
		})

		is.Equal(len(result1), 1)
		is.Equal(len(result2), 4)
		is.Equal(result2, map[string]int{"1": 1, "2": 2, "3": 3, "4": 4})
	}

	{
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
}

// TestMapValues Transforms a map to other type map by manipulating it's values.
// func MapValues[K comparable, V any, T any](m map[K]V, iteratee func(key K, value V) T) map[K]T
func TestMapValues(t *testing.T) {
	{
		t.Parallel()
		is := assert.New(t)

		result1 := lo.MapValues(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
			return "Hello"
		})
		result2 := lo.MapValues(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
			return strconv.FormatInt(int64(x), 10)
		})

		is.Equal(len(result1), 4)
		is.Equal(len(result2), 4)
		is.Equal(result1, map[int]string{1: "Hello", 2: "Hello", 3: "Hello", 4: "Hello"})
		is.Equal(result2, map[int]string{1: "1", 2: "2", 3: "3", 4: "4"})
	}

	{
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
	t.Parallel()
	is := assert.New(t)

	{
		r1 := lo.Entries(map[string]int{"foo": 1, "bar": 2})

		sort.Slice(r1, func(i, j int) bool {
			return r1[i].Value < r1[j].Value
		})
		is.EqualValues(r1, []lo.Entry[string, int]{
			{
				Key:   "foo",
				Value: 1,
			},
			{
				Key:   "bar",
				Value: 2,
			},
		})
	}

	{
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
}

func mapEntriesTest[I any, O any](t *testing.T, in map[string]I, iteratee func(string, I) (string, O), expected map[string]O) {
	is := assert.New(t)
	result := lo.MapEntries(in, iteratee)
	is.Equal(result, expected)
}

func TestMapEntries(t *testing.T) {
	mapEntriesTest(t, map[string]int{"foo": 1, "bar": 2}, func(k string, v int) (string, int) {
		return k, v + 1
	}, map[string]int{"foo": 2, "bar": 3})
	mapEntriesTest(t, map[string]int{"foo": 1, "bar": 2}, func(k string, v int) (string, string) {
		return k, k + strconv.Itoa(v)
	}, map[string]string{"foo": "foo1", "bar": "bar2"})
	mapEntriesTest(t, map[string]int{"foo": 1, "bar": 2}, func(k string, v int) (string, string) {
		return k, strconv.Itoa(v) + k
	}, map[string]string{"foo": "1foo", "bar": "2bar"})

	// NoMutation
	{
		is := assert.New(t)
		r1 := map[string]int{"foo": 1, "bar": 2}
		lo.MapEntries(r1, func(k string, v int) (string, string) {
			return k, strconv.Itoa(v) + "!!"
		})
		is.Equal(r1, map[string]int{"foo": 1, "bar": 2})
	}
	// EmptyInput
	{
		mapEntriesTest(t, map[string]int{}, func(k string, v int) (string, string) {
			return k, strconv.Itoa(v) + "!!"
		}, map[string]string{})

		mapEntriesTest(t, map[string]any{}, func(k string, v any) (string, any) {
			return k, v
		}, map[string]any{})
	}
	// Identity
	{
		mapEntriesTest(t, map[string]int{"foo": 1, "bar": 2}, func(k string, v int) (string, int) {
			return k, v
		}, map[string]int{"foo": 1, "bar": 2})
		mapEntriesTest(t, map[string]any{"foo": 1, "bar": "2", "ccc": true}, func(k string, v any) (string, any) {
			return k, v
		}, map[string]any{"foo": 1, "bar": "2", "ccc": true})
	}
	// ToConstantEntry
	{
		mapEntriesTest(t, map[string]any{"foo": 1, "bar": "2", "ccc": true}, func(k string, v any) (string, any) {
			return "key", "value"
		}, map[string]any{"key": "value"})
		mapEntriesTest(t, map[string]any{"foo": 1, "bar": "2", "ccc": true}, func(k string, v any) (string, any) {
			return "b", 5
		}, map[string]any{"b": 5})
	}

	//// OverlappingKeys
	//// because using range over map, the order is not guaranteed
	//// this test is not deterministic
	//{
	//	mapEntriesTest(t, map[string]any{"foo": 1, "foo2": 2, "Foo": 2, "Foo2": "2", "bar": "2", "ccc": true}, func(k string, v any) (string, any) {
	//		return string(k[0]), v
	//	}, map[string]any{"F": "2", "b": "2", "c": true, "f": 2})
	//	mapEntriesTest(t, map[string]string{"foo": "1", "foo2": "2", "Foo": "2", "Foo2": "2", "bar": "2", "ccc": "true"}, func(k string, v string) (string, string) {
	//		return v, k
	//	}, map[string]string{"1": "foo", "2": "bar", "true": "ccc"})
	//}
	//NormalMappers
	{
		mapEntriesTest(t, map[string]string{"foo": "1", "foo2": "2", "Foo": "2", "Foo2": "2", "bar": "2", "ccc": "true"}, func(k string, v string) (string, string) {
			return k, k + v
		}, map[string]string{"Foo": "Foo2", "Foo2": "Foo22", "bar": "bar2", "ccc": "ccctrue", "foo": "foo1", "foo2": "foo22"})

		mapEntriesTest(t, map[string]struct {
			name string
			age  int
		}{"1-11-1": {name: "foo", age: 1}, "2-22-2": {name: "bar", age: 2}}, func(k string, v struct {
			name string
			age  int
		},
		) (string, string) {
			return v.name, k
		}, map[string]string{"bar": "2-22-2", "foo": "1-11-1"})
	}
}

func TestToPairs(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.ToPairs(map[string]int{"baz": 3, "qux": 4})

	sort.Slice(r1, func(i, j int) bool {
		return r1[i].Value < r1[j].Value
	})
	is.EqualValues(r1, []lo.Entry[string, int]{
		{
			Key:   "baz",
			Value: 3,
		},
		{
			Key:   "qux",
			Value: 4,
		},
	})
}
