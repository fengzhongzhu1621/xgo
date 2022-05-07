package utils

import (
	"reflect"
	"testing"
)

func TestCopyAndInsensitiviseMap(t *testing.T) {
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

func TestFlattenAndMergeMap(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
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

		expect = map[string]interface{}{
			"a": 32,
			"b": map[string]interface{}{
				"c": "A",
				"d": map[string]interface{}{
					"e": "E",
					"f": "F",
				},
			},
		}
		if !reflect.DeepEqual(m, expect) {
			t.Fatal("DeepSearch error")
		}
	})

	t.Run("EmptyMap", func(t *testing.T) {
		m := map[string]interface{}{}
		lastMap := DeepSearch(m, []string{"b", "d"})
		// 验证最后一个字典的值
		expect := map[string]interface{}{}
		if !reflect.DeepEqual(lastMap, expect) {
			t.Fatalf("DeepSearch error actual is %v, expect is %v", lastMap, expect)
		}
		// 验证生成后的深度遍历后的字典
		expect = map[string]interface{}{
			"b": map[string]interface{}{
				"d": map[string]interface{}{},
			},
		}
		if !reflect.DeepEqual(m, expect) {
			t.Fatal("DeepSearch error")
		}
	})
}

func TestMergeFlatMap(t *testing.T) {
	t.Parallel()

	t.Run("UpperKey", func(t *testing.T) {
		var shadow = map[string]bool{}
		var src = map[string]interface{}{
			"A":     1,
			"A_B":   2,
			"A_B_C": 3,
			"a":     4,
		}
		keyDelim := "_"
		actual := MergeFlatMap(shadow, src, keyDelim)
		expect := map[string]bool{
			"a":     true,
			"a_b":   true,
			"a_b_c": true,
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Fatalf("MergeFlatMap error actual is %v, expect is %v", actual, expect)
		}
	})

	t.Run("LowerKey", func(t *testing.T) {
		var shadow = map[string]bool{}
		var src = map[string]interface{}{
			"a":     1,
			"a_b":   2,
			"a_b_C": 3,
			"A":     4,
		}
		keyDelim := "_"
		actual := MergeFlatMap(shadow, src, keyDelim)
		expect := map[string]bool{
			"a": true,
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Fatalf("MergeFlatMap error actual is %v, expect is %v", actual, expect)
		}
	})
}

func TestMergeMaps(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		var src = map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": map[string]interface{}{
				"c1": 3,
			},
			"d": map[interface{}]interface{}{
				"d1": 4,
			},
		}
		dst := map[string]interface{}{
			"A": 11,
			"b": 22,
			"C": map[string]interface{}{
				"c1": 33,
			},
		}
		dst2 := map[interface{}]interface{}{}

		MergeMaps(src, dst, dst2)

		expectDst := map[string]interface{}{
			"A": 1,
			"b": 2,
			"C": map[string]interface{}{
				"c1": 3,
			},
			"d": map[interface{}]interface{}{
				"d1": 4,
			},
		}
		if !reflect.DeepEqual(dst, expectDst) {
			t.Fatalf("MergeMaps error actual is %v, expect is %v", dst, expectDst)
		}

		expectDst2 := map[interface{}]interface{}{
			"A": 1,
			"b": 2,
			"d": map[interface{}]interface{}{
				"d1": 4,
			},
		}
		if !reflect.DeepEqual(dst2, expectDst2) {
			t.Fatalf("MergeMaps error actual is %v, expect is %v", dst2, expectDst2)
		}
	})
}

func TestCreateDeepMap(t *testing.T) {
	t.Run("UpperKey", func(t *testing.T) {
		var src = map[string]interface{}{
			"A":     1,
			"A_B":   2,
			"A_B_C": 3,
			"d":     4,
		}
		keyDelim := "_"
		actual := CreateDeepMap(src, keyDelim)
		expect := map[string]interface{}{
			"A": map[string]interface{}{
				"B": map[string]interface{}{
					"c": 3,
				},
				"b": 2,
			},
			"a": 1,
			"d": 4,
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Fatalf("CreateDeepMap error actual is %v, expect is %v", actual, expect)
		}
	})
}

func TestSetDeepMapValue(t *testing.T) {
	t.Parallel()

	t.Run("KeyNotExists", func(t *testing.T) {
		var src = map[string]interface{}{
			"A": 1,
			"b": 2,
		}
		keyDelim := "_"
		key := "A"
		value := 11
		SetDeepMapValue(src, key, value, keyDelim)
		expect := map[string]interface{}{
			"A": 1,
			"a": 11,
			"b": 2,
		}
		if !reflect.DeepEqual(src, expect) {
			t.Fatalf("SetDeepMapValue error actual is %v, expect is %v", src, expect)
		}
	})

	t.Run("KeyExists", func(t *testing.T) {
		var src = map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{},
			},
		}
		keyDelim := "_"
		key := "a_b_c"
		value := 3
		SetDeepMapValue(src, key, value, keyDelim)
		expect := map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": 3,
				},
			},
		}
		if !reflect.DeepEqual(src, expect) {
			t.Fatalf("SetDeepMapValue error actual is %v, expect is %v", src, expect)
		}
	})
}
