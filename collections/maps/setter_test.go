package maps

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SetItems 成功的情况
func TestSetItemsSuccessCase(t *testing.T) {
	// depth 1，val type int
	err := SetItems(deploySpec, "intKey4SetItem", 5)
	assert.Nil(t, err)
	ret, _ := GetItems(deploySpec, []string{"intKey4SetItem"})
	assert.Equal(t, 5, ret)

	// depth 2, val type string
	err = SetItems(deploySpec, "strategy.type", "Rolling")
	assert.Nil(t, err)
	ret, _ = GetItems(deploySpec, []string{"strategy", "type"})
	assert.Equal(t, "Rolling", ret)

	// depth 3, val type string
	err = SetItems(deploySpec, []string{"template", "spec", "restartPolicy"}, "Never")
	assert.Nil(t, err)
	ret, _ = GetItems(deploySpec, []string{"template", "spec", "restartPolicy"})
	assert.Equal(t, "Never", ret)

	// key noy exists
	err = SetItems(deploySpec, []string{"selector", "testKey"}, "testVal")
	assert.Nil(t, err)
	ret, _ = GetItems(deploySpec, "selector.testKey")
	assert.Equal(t, "testVal", ret)
}

// SetItems 失败的情况
func TestSetItemsFailCase(t *testing.T) {
	// invalid paths type error
	err := SetItems(deploySpec, 0, 1)
	assert.True(t, errors.Is(err, ErrInvalidPathType))

	// not paths error
	err = SetItems(deploySpec, []string{}, 1)
	assert.NotNil(t, err)

	// not map[string]interface{} type error
	err = SetItems(deploySpec, []string{"replicas", "testKey"}, 1)
	assert.NotNil(t, err)

	// key not exist
	err = SetItems(deploySpec, []string{"templateKey", "spec"}, 1)
	assert.NotNil(t, err)

	err = SetItems(deploySpec, "templateKey.spec", 123)
	assert.NotNil(t, err)

	// paths type error
	err = SetItems(deploySpec, []int{123, 456}, 1)
	assert.NotNil(t, err)

	err = SetItems(deploySpec, 123, 1)
	assert.NotNil(t, err)
}

func TestSetDeepMapValue(t *testing.T) {
	t.Parallel()

	t.Run("KeyNotExists", func(t *testing.T) {
		src := map[string]interface{}{
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
		src := map[string]interface{}{
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
