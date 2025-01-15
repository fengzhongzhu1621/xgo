package maps

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestRangeKey 遍历字典的 key
func TestRangeKey(t *testing.T) {
	// 创建一个 map
	myMap := map[string]int{
		"apple":  5,
		"banana": 7,
		"orange": 3,
	}

	// 遍历 map 的 key
	for key := range myMap {
		fmt.Println("Key:", key)
	}
}

// TestKeys Returns a slice of the map's keys.
// func Keys[K comparable, V any](m map[K]V) []K
func TestKeys(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "a",
		3: "b",
		4: "c",
		5: "d",
	}

	keys := maputil.Keys(m)
	sort.Ints(keys)

	assert.Equal(t, keys, []int{1, 2, 3, 4, 5})
}

// TestKeysBy Creates a slice whose element is the result of function mapper invoked by every map's key.
// func KeysBy[K comparable, V any, T any](m map[K]V, mapper func(item K) T) []T
func TestKeysBy(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "a",
		3: "b",
	}

	keys := maputil.KeysBy(m, func(n int) int {
		return n + 1
	})

	sort.Ints(keys)

	fmt.Println(keys)

	// Output:
	// [2 3 4]
}

func TestKeys2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.Keys(map[string]int{"foo": 1, "bar": 2})
	sort.Strings(r1)
	is.Equal(r1, []string{"bar", "foo"})

	r2 := lo.Keys(map[string]int{})
	is.Empty(r2)

	r3 := lo.Keys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
	sort.Strings(r3)
	is.Equal(r3, []string{"bar", "baz", "foo"})

	r4 := lo.Keys[string, int]()
	is.Equal(r4, []string{})

	r5 := lo.Keys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"bar": 3})
	sort.Strings(r5)
	is.Equal(r5, []string{"bar", "bar", "foo"})
}

func TestUniqKeys(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.UniqKeys(map[string]int{"foo": 1, "bar": 2})
	sort.Strings(r1)
	is.Equal(r1, []string{"bar", "foo"})

	r2 := lo.UniqKeys(map[string]int{})
	is.Empty(r2)

	r3 := lo.UniqKeys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
	sort.Strings(r3)
	is.Equal(r3, []string{"bar", "baz", "foo"})

	r4 := lo.UniqKeys[string, int]()
	is.Equal(r4, []string{})

	r5 := lo.UniqKeys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"foo": 1, "bar": 3})
	sort.Strings(r5)
	is.Equal(r5, []string{"bar", "foo"})

	// check order
	r6 := lo.UniqKeys(map[string]int{"foo": 1}, map[string]int{"bar": 3})
	is.Equal(r6, []string{"foo", "bar"})
}

// TestToCaseInsensitiveValue 将字典的key转换为小写，返回新的字典
func TestToCaseInsensitiveValue(t *testing.T) {
	t.Parallel()

	var (
		given = map[string]interface{}{
			"Foo": 32,
			"Bar": map[interface{}]interface{}{
				"ABc": "A",
				"cDE": "B",
			},
		}
		expected = map[string]interface{}{
			"foo": 32,
			"bar": map[string]interface{}{
				"abc": "A",
				"cde": "B",
			},
		}
	)

	got := ToCaseInsensitiveValue(given)

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Got %q\nexpected\n%q", got, expected)
	}

	if _, ok := given["foo"]; ok {
		t.Fatal("Input map changed")
	}

	if _, ok := given["bar"]; ok {
		t.Fatal("Input map changed")
	}

	m := given["Bar"].(map[interface{}]interface{})
	if _, ok := m["ABc"]; !ok {
		t.Fatal("Input map changed")
	}
}

// TestHasKey Checks if map has key or not. This function is used to replace the following boilerplate code
// func HasKey[K comparable, V any](m map[K]V, key K) bool
func TestHasKey(t *testing.T) {
	{
		obj := map[string]interface{}{
			"key_1": "val_1",
			"key_3": "val_3",
		}
		assert.True(t, ExistsKey(obj, "key_1"))
		assert.False(t, ExistsKey(obj, "key_2"))
		assert.True(t, ExistsKey(obj, "key_3"))
	}

	{
		m := map[string]int{
			"a": 1,
			"b": 2,
		}

		result1 := maputil.HasKey(m, "a")
		result2 := maputil.HasKey(m, "c")

		fmt.Println(result1)
		fmt.Println(result2)

		// Output:
		// true
		// false
	}

	{
		t.Parallel()
		is := assert.New(t)

		r1 := lo.HasKey(map[string]int{"foo": 1}, "bar")
		is.False(r1)

		r2 := lo.HasKey(map[string]int{"foo": 1}, "foo")
		is.True(r2)
	}

}

// TestMapGetOrSet Returns value of the given key or set the given value value if not present.
// func GetOrSet[K comparable, V any](m map[K]V, key K, value V) V
func TestMapGetOrSet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{

		r1 := lo.ValueOr(map[string]int{"foo": 1}, "bar", 2)
		is.Equal(r1, 2)

		r2 := lo.ValueOr(map[string]int{"foo": 1}, "foo", 2)
		is.Equal(r2, 1)
	}

	{
		m := map[int]string{
			1: "a",
		}

		result1 := maputil.GetOrSet(m, 1, "1")
		result2 := maputil.GetOrSet(m, 2, "b")

		fmt.Println(result1)
		fmt.Println(result2)

		// Output:
		// a
		// b
	}
}
