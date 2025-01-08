package maps

import (
	"reflect"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/stretchr/testify/assert"
)

// TestMerge 合并多个 maps 相同的 key 会被后来的 key 覆盖
func TestMerge(t *testing.T) {
	m1 := map[int]string{
		1: "a",
		2: "b",
	}
	m2 := map[int]string{
		1: "c",
		3: "d",
	}

	result := maputil.Merge(m1, m2)

	assert.Equal(t, result, map[int]string{
		1: "c",
		2: "b",
		3: "d",
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

func TestMergeMaps(t *testing.T) {
	t.Parallel()

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
