package maps

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/stretchr/testify/assert"
)

// TestExistsKey 判断 key 是否存在
func TestExistsKey(t *testing.T) {
	obj := map[string]interface{}{
		"key_1": "val_1",
		"key_3": "val_3",
	}
	assert.True(t, ExistsKey(obj, "key_1"))
	assert.False(t, ExistsKey(obj, "key_2"))
	assert.True(t, ExistsKey(obj, "key_3"))
}

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
