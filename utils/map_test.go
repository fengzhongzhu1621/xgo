package utils

import (
	"reflect"
	"testing"
)

func TestCopyAndInsensitiviseMap(t *testing.T) {
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

func TestFlattenAndMergeMap(t *testing.T) {
	var data = map[string]interface{}{
		"KEY": map[string]interface{}{
			"a": 1,
			"b": 2,
		},
	}
	actual := FlattenAndMergeMap(nil, data, "", "_")
	expect := map[string]interface{}{
		"key_a": 1,
		"key_b": 2,
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Fatal("FlattenAndMergeMap error")
	}
}

func TestDeepSearch(t *testing.T) {
	m := map[string]interface{}{
		"a": 32,
		"b": map[string]interface{}{
			"c": "A",
			"d": map[string]interface{}{
				"e": "E",
				"f": "F",
			},
		},
	}
	actual := DeepSearch(m, []string{"b", "d"})
	expect := map[string]interface{}{
		"e": "E",
		"f": "F",
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Fatal("DeepSearch error")
	}
}
