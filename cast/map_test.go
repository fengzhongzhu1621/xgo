package cast

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/util/gconv"

	"github.com/duke-git/lancet/v2/structs"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gookit/goutil/arrutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestStr2map(t *testing.T) {
	s := "a=1&b=2&c="
	actual := Str2map(s, "&", "=")
	expect := map[string]string{"a": "1", "b": "2", "c": ""}
	assert.Equal(t, expect, actual)
}

// TestKeyBy Converts a slice to a map based on a callback function.
// func KeyBy[T any, U comparable](slice []T, iteratee func(item T) U) map[U]T
func TestSliceToMapKeyBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	{
		result1 := lo.KeyBy([]string{"a", "aa", "aaa"}, func(str string) int {
			return len(str)
		})
		is.Equal(result1, map[int]string{1: "a", 2: "aa", 3: "aaa"})
	}

	{
		result := slice.KeyBy([]string{"a", "ab", "abc"}, func(str string) int {
			return len(str)
		})
		assert.Equal(t, map[int]string{1: "a", 2: "ab", 3: "abc"}, result)
	}

}

func TestStructToMap(t *testing.T) {
	{
		// TestToMap convert a valid struct to a map
		// func (s *Struct) ToMap() (map[string]any, error)
		// func ToMap(v any) (map[string]any, error)
		type People struct {
			Name string `json:"name"`
		}
		p1 := &People{Name: "11"}
		// use constructor function
		s1 := structs.New(p1)
		m1, _ := s1.ToMap()
		fmt.Println(m1) // map[name:11]

		// use static function
		m2, _ := structs.ToMap(p1)
		fmt.Println(m2) // map[name:11]
	}

	{
		type People struct {
			Name string `json:"name"`
			age  int
		}
		p := People{
			"test",
			100,
		}
		pm, _ := convertor.StructToMap(p)

		fmt.Println(pm)

		// Output:
		// map[name:test]
	}

	{
		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		user := User{Name: "John", Age: 28}

		// 结构体转换成 map
		userMap := gconv.Map(user)
		fmt.Println(userMap) // 输出：map[age:28 name:John]
	}
}

// func StructToMap(value any) (map[string]any, error)
func TestStructSliceToMap(t *testing.T) {
	{
		// TestToMap 结构体转换为 MAP
		// Convert a slice of structs to a map based on iteratee function.
		// func ToMap[T any, K comparable, V any](array []T, iteratee func(T) (K, V)) map[K]V
		type Message struct {
			name string
			code int
		}
		messages := []Message{
			{name: "Hello", code: 100},
			{name: "Hi", code: 101},
		}

		result := convertor.ToMap(messages, func(msg Message) (int, string) {
			return msg.code, msg.name
		})

		fmt.Println(result)

		// Output:
		// map[100:Hello 101:Hi]
	}

	t.Parallel()
	{
		type foo struct {
			baz string
			bar int
		}
		transform := func(f *foo) (string, int) {
			return f.baz, f.bar
		}
		testCases := []struct {
			in     []*foo
			expect map[string]int
		}{
			{
				in:     []*foo{{baz: "apple", bar: 1}},
				expect: map[string]int{"apple": 1},
			},
			{
				in:     []*foo{{baz: "apple", bar: 1}, {baz: "banana", bar: 2}},
				expect: map[string]int{"apple": 1, "banana": 2},
			},
			{
				in:     []*foo{{baz: "apple", bar: 1}, {baz: "apple", bar: 2}},
				expect: map[string]int{"apple": 2},
			},
		}
		for i, testCase := range testCases {
			t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
				is := assert.New(t)
				is.Equal(lo.Associate(testCase.in, transform), testCase.expect)
			})
		}
	}
}

// 遍历结构体数组，转换为指定格式的map
func TestSliceToMap(t *testing.T) {
	t.Parallel()

	type foo struct {
		baz string
		bar int
	}
	transform := func(f *foo) (string, int) {
		return f.baz, f.bar
	}
	testCases := []struct {
		in     []*foo
		expect map[string]int
	}{
		{
			in:     []*foo{{baz: "apple", bar: 1}},
			expect: map[string]int{"apple": 1},
		},
		{
			in:     []*foo{{baz: "apple", bar: 1}, {baz: "banana", bar: 2}},
			expect: map[string]int{"apple": 1, "banana": 2},
		},
		{
			in:     []*foo{{baz: "apple", bar: 1}, {baz: "apple", bar: 2}},
			expect: map[string]int{"apple": 2},
		},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			is := assert.New(t)
			is.Equal(lo.SliceToMap(testCase.in, transform), testCase.expect)
		})
	}
}

// TestFrequency Counts the frequency of each element in the slice.
// func Frequency[T comparable](slice []T) map[T]int
func TestFrequency(t *testing.T) {
	{
		strs := []string{"a", "b", "b", "c", "c", "c"}
		result := slice.Frequency(strs)

		assert.Equal(t, map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}, result)
	}

	{
		t.Parallel()
		is := assert.New(t)

		is.Equal(map[int]int{}, lo.CountValues([]int{}))
		is.Equal(map[int]int{1: 1, 2: 1}, lo.CountValues([]int{1, 2}))
		is.Equal(map[int]int{1: 1, 2: 2}, lo.CountValues([]int{1, 2, 2}))
		is.Equal(map[string]int{"": 1, "foo": 1, "bar": 1}, lo.CountValues([]string{"foo", "bar", ""}))
		is.Equal(map[string]int{"foo": 1, "bar": 2}, lo.CountValues([]string{"foo", "bar", "bar"}))
	}
}

func TestCountValuesBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	oddEven := func(v int) bool {
		return v%2 == 0
	}
	length := func(v string) int {
		return len(v)
	}

	result1 := lo.CountValuesBy([]int{}, oddEven)
	result2 := lo.CountValuesBy([]int{1, 2}, oddEven)
	result3 := lo.CountValuesBy([]int{1, 2, 2}, oddEven)
	result4 := lo.CountValuesBy([]string{"foo", "bar", ""}, length)
	result5 := lo.CountValuesBy([]string{"foo", "bar", "bar"}, length)

	is.Equal(map[bool]int{}, result1)
	is.Equal(map[bool]int{true: 1, false: 1}, result2)
	is.Equal(map[bool]int{true: 2, false: 1}, result3)
	is.Equal(map[int]int{0: 1, 3: 2}, result4)
	is.Equal(map[int]int{3: 3}, result5)
}

func TestZip(t *testing.T) {
	keys := []string{"key0", "key1"}

	{
		mp := arrutil.CombineToMap(keys, []int{1, 2})
		assert.Len(t, mp, 2)
		assert.Equal(t, 1, mp["key0"])
		assert.Equal(t, 2, mp["key1"])

		mp = arrutil.CombineToMap(keys, []int{1})
		assert.Len(t, mp, 1)
		assert.Equal(t, 1, mp["key0"])
	}

	{
		mp := arrutil.CombineToSMap(keys, []string{"val0", "val1"})
		assert.Len(t, mp, 2)
		assert.Equal(t, "val0", mp["key0"])

		mp = arrutil.CombineToSMap(keys, []string{"val0"})
		assert.Len(t, mp, 2)
		assert.Equal(t, "val0", mp["key0"])
		assert.Equal(t, "", mp["key1"])
	}
}

func TestStringMapStringSliceE(t *testing.T) {
	// ToStringMapString inputs/outputs
	var stringMapString = map[string]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var stringMapInterface = map[string]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapString = map[interface{}]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapInterface = map[interface{}]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}

	// ToStringMapStringSlice inputs/outputs
	var stringMapStringSlice = map[string][]string{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var stringMapInterfaceSlice = map[string][]interface{}{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var stringMapInterfaceInterfaceSlice = map[string]interface{}{"key 1": []interface{}{"value 1", "value 2", "value 3"}, "key 2": []interface{}{"value 1", "value 2", "value 3"}, "key 3": []interface{}{"value 1", "value 2", "value 3"}}
	var stringMapStringSingleSliceFieldsResult = map[string][]string{"key 1": {"value", "1"}, "key 2": {"value", "2"}, "key 3": {"value", "3"}}
	var interfaceMapStringSlice = map[interface{}][]string{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var interfaceMapInterfaceSlice = map[interface{}][]interface{}{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}

	var stringMapStringSliceMultiple = map[string][]string{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var stringMapStringSliceSingle = map[string][]string{"key 1": {"value 1"}, "key 2": {"value 2"}, "key 3": {"value 3"}}

	var stringMapInterface1 = map[string]interface{}{"key 1": []string{"value 1"}, "key 2": []string{"value 2"}}
	var stringMapInterfaceResult1 = map[string][]string{"key 1": {"value 1"}, "key 2": {"value 2"}}

	var jsonStringMapString = `{"key 1": "value 1", "key 2": "value 2"}`
	var jsonStringMapStringArray = `{"key 1": ["value 1"], "key 2": ["value 2", "value 3"]}`
	var jsonStringMapStringArrayResult = map[string][]string{"key 1": {"value 1"}, "key 2": {"value 2", "value 3"}}

	type Key struct {
		k string
	}

	tests := []struct {
		input  interface{}
		expect map[string][]string
		iserr  bool
	}{
		{stringMapStringSlice, stringMapStringSlice, false},
		{stringMapInterfaceSlice, stringMapStringSlice, false},
		{stringMapInterfaceInterfaceSlice, stringMapStringSlice, false},
		{stringMapStringSliceMultiple, stringMapStringSlice, false},
		{stringMapStringSliceMultiple, stringMapStringSlice, false},
		{stringMapString, stringMapStringSliceSingle, false},
		{stringMapInterface, stringMapStringSliceSingle, false},
		{stringMapInterface1, stringMapInterfaceResult1, false},
		{interfaceMapStringSlice, stringMapStringSlice, false},
		{interfaceMapInterfaceSlice, stringMapStringSlice, false},
		{interfaceMapString, stringMapStringSingleSliceFieldsResult, false},
		{interfaceMapInterface, stringMapStringSingleSliceFieldsResult, false},
		{jsonStringMapStringArray, jsonStringMapStringArrayResult, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{map[interface{}]interface{}{"foo": testing.T{}}, nil, true},
		{map[interface{}]interface{}{Key{"foo"}: "bar"}, nil, true}, // ToStringE(Key{"foo"}) should fail
		{jsonStringMapString, nil, true},
		{"", nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringMapStringSliceE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToStringMapStringSlice(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToStringMapE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect map[string]interface{}
		iserr  bool
	}{
		{map[interface{}]interface{}{"tag": "tags", "group": "groups"}, map[string]interface{}{"tag": "tags", "group": "groups"}, false},
		{map[string]interface{}{"tag": "tags", "group": "groups"}, map[string]interface{}{"tag": "tags", "group": "groups"}, false},
		{`{"tag": "tags", "group": "groups"}`, map[string]interface{}{"tag": "tags", "group": "groups"}, false},
		{`{"tag": "tags", "group": true}`, map[string]interface{}{"tag": "tags", "group": true}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{"", nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringMapE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToStringMap(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToStringMapBoolE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect map[string]bool
		iserr  bool
	}{
		{map[interface{}]interface{}{"v1": true, "v2": false}, map[string]bool{"v1": true, "v2": false}, false},
		{map[string]interface{}{"v1": true, "v2": false}, map[string]bool{"v1": true, "v2": false}, false},
		{map[string]bool{"v1": true, "v2": false}, map[string]bool{"v1": true, "v2": false}, false},
		{`{"v1": true, "v2": false}`, map[string]bool{"v1": true, "v2": false}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{"", nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringMapBoolE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToStringMapBool(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToStringMapIntE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect map[string]int
		iserr  bool
	}{
		{map[interface{}]interface{}{"v1": 1, "v2": 222}, map[string]int{"v1": 1, "v2": 222}, false},
		{map[string]interface{}{"v1": 342, "v2": 5141}, map[string]int{"v1": 342, "v2": 5141}, false},
		{map[string]int{"v1": 33, "v2": 88}, map[string]int{"v1": 33, "v2": 88}, false},
		{map[string]int32{"v1": int32(33), "v2": int32(88)}, map[string]int{"v1": 33, "v2": 88}, false},
		{map[string]uint16{"v1": uint16(33), "v2": uint16(88)}, map[string]int{"v1": 33, "v2": 88}, false},
		{map[string]float64{"v1": float64(8.22), "v2": float64(43.32)}, map[string]int{"v1": 8, "v2": 43}, false},
		{`{"v1": 67, "v2": 56}`, map[string]int{"v1": 67, "v2": 56}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{"", nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringMapIntE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToStringMapInt(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToStringMapInt64E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect map[string]int64
		iserr  bool
	}{
		{map[interface{}]interface{}{"v1": int32(8), "v2": int32(888)}, map[string]int64{"v1": int64(8), "v2": int64(888)}, false},
		{map[string]interface{}{"v1": int64(45), "v2": int64(67)}, map[string]int64{"v1": 45, "v2": 67}, false},
		{map[string]int64{"v1": 33, "v2": 88}, map[string]int64{"v1": 33, "v2": 88}, false},
		{map[string]int{"v1": 33, "v2": 88}, map[string]int64{"v1": 33, "v2": 88}, false},
		{map[string]int32{"v1": int32(33), "v2": int32(88)}, map[string]int64{"v1": 33, "v2": 88}, false},
		{map[string]uint16{"v1": uint16(33), "v2": uint16(88)}, map[string]int64{"v1": 33, "v2": 88}, false},
		{map[string]float64{"v1": float64(8.22), "v2": float64(43.32)}, map[string]int64{"v1": 8, "v2": 43}, false},
		{`{"v1": 67, "v2": 56}`, map[string]int64{"v1": 67, "v2": 56}, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{"", nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringMapInt64E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToStringMapInt64(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToStringMapStringE(t *testing.T) {
	var stringMapString = map[string]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var stringMapInterface = map[string]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapString = map[interface{}]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapInterface = map[interface{}]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var jsonString = `{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}`
	var invalidJsonString = `{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"`
	var emptyString = ""

	tests := []struct {
		input  interface{}
		expect map[string]string
		iserr  bool
	}{
		{stringMapString, stringMapString, false},
		{stringMapInterface, stringMapString, false},
		{interfaceMapString, stringMapString, false},
		{interfaceMapInterface, stringMapString, false},
		{jsonString, stringMapString, false},

		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{invalidJsonString, nil, true},
		{emptyString, nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringMapStringE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToStringMapString(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}
