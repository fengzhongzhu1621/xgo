package maps

import (
	"reflect"
	"testing"
)

func TestCreateDeepMap(t *testing.T) {
	t.Parallel()

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
