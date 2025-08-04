package maps

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeepMap(t *testing.T) {
	t.Parallel()

	t.Run("UpperKey", func(t *testing.T) {
		src := map[string]interface{}{
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

// TestFromEntries Creates a map based on a slice of key/value pairs.
//
//	type Entry[K comparable, V any] struct {
//	    Key   K
//	    Value V
//	}
//
// func FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V
func TestFromEntries(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		r1 := lo.FromEntries([]lo.Entry[string, int]{
			{
				Key:   "foo",
				Value: 1,
			},
			{
				Key:   "bar",
				Value: 2,
			},
		})

		is.Len(r1, 2)
		is.Equal(r1["foo"], 1)
		is.Equal(r1["bar"], 2)
	}

	{
		result := maputil.FromEntries([]maputil.Entry[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
			{Key: "c", Value: 3},
		})

		fmt.Println(result)

		// Output:
		// map[a:1 b:2 c:3]
	}
}

func TestFromPairs(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.FromPairs([]lo.Entry[string, int]{
		{
			Key:   "baz",
			Value: 3,
		},
		{
			Key:   "qux",
			Value: 4,
		},
	})

	is.Len(r1, 2)
	is.Equal(r1["baz"], 3)
	is.Equal(r1["qux"], 4)
}
