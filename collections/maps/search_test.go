package maps

import (
	"reflect"
	"testing"
)

func TestSearchMap(t *testing.T) {
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
		actual := SearchMap(m, []string{"b", "d"})
		expect := map[string]interface{}{
			"e": "E",
			"f": "F",
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Fatal("DeepSearch error")
		}

		actual = SearchMap(m, []string{"b", "d", "e"})
		expect2 := "E"
		if !reflect.DeepEqual(actual, expect2) {
			t.Fatal("DeepSearch error")
		}
	})

	t.Run("NoValue", func(t *testing.T) {
		m := map[string]interface{}{
			"a": 32,
		}
		actual := SearchMap(m, []string{"b", "d", "e"})
		if !reflect.DeepEqual(actual, nil) {
			t.Fatal("DeepSearch error")
		}
	})
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

func TestSearchIndexableWithPathPrefixes(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		src := map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{
				"c": "3",
				"d": []interface{}{
					"5",
					"6",
				},
			},
		}
		keyDelim := "_"
		path := []string{"b", "d", "1"}
		actual := SearchIndexableWithPathPrefixes(src, path, keyDelim)
		expect := "6"
		if !reflect.DeepEqual(actual, expect) {
			t.Fatalf(
				"SearchIndexableWithPathPrefixes error actual is %v, expect is %v",
				actual,
				expect,
			)
		}
	})
}

func TestIsPathShadowedInDeepMap(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		src := map[string]interface{}{
			"a": "1",
			"b": map[string]interface{}{
				"c": "3",
				"d": map[string]interface{}{
					"e": "5",
					"f": "6",
				},
			},
		}
		keyDelim := "_"
		path := []string{"b", "c", "c1", "c2"}
		actual := IsPathShadowedInDeepMap(path, src, keyDelim)
		expect := "b_c"
		if !reflect.DeepEqual(actual, expect) {
			t.Fatalf("IsPathShadowedInDeepMap error actual is %v, expect is %v", actual, expect)
		}
	})
}

func TestIsPathShadowedInFlatMap(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		src := map[string]interface{}{
			"B":   1,
			"A_B": 2,
		}
		keyDelim := "_"
		path := []string{"A", "B", "C"}
		actual := IsPathShadowedInFlatMap(path, src, keyDelim)
		expect := "A_B"
		if !reflect.DeepEqual(actual, expect) {
			t.Fatalf("IsPathShadowedInFlatMap error actual is %v, expect is %v", actual, expect)
		}
	})
}
