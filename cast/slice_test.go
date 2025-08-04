package cast

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/gookit/goutil/arrutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestChannelToSlice(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 10*time.Millisecond)
	is := assert.New(t)

	ch := lo.SliceToChannel(2, []int{1, 2, 3})
	items := lo.ChannelToSlice(ch)

	is.Equal([]int{1, 2, 3}, items)
}

// func MapToSlice[T any, K comparable, V any](aMap map[K]V, iteratee func(K, V) T) []T
func TestMapToSlice(t *testing.T) {
	aMap := map[string]int{"a": 1, "b": 2, "c": 3}
	result := convertor.MapToSlice(aMap, func(key string, value int) string {
		return key + ":" + strconv.Itoa(value)
	})

	fmt.Println(result) //[]string{"a:1", "b:2", "c:3"}
}

func TestToAnySlice(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	in1 := []int{0, 1, 2, 3}
	in2 := []int{}
	out1 := lo.ToAnySlice(in1)
	out2 := lo.ToAnySlice(in2)

	is.Equal([]any{0, 1, 2, 3}, out1)
	is.Equal([]any{}, out2)
}

func TestFromAnySlice(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.NotPanics(func() {
		out1, ok1 := lo.FromAnySlice[string]([]any{"foobar", 42})
		out2, ok2 := lo.FromAnySlice[string]([]any{"foobar", "42"})

		is.Equal([]string{}, out1)
		is.False(ok1)
		is.Equal([]string{"foobar", "42"}, out2)
		is.True(ok2)
	})
}

func TestStringsToInts(t *testing.T) {
	is := assert.New(t)

	ints, err := arrutil.StringsToInts([]string{"1", "2"})
	is.Nil(err)
	is.Equal("[]int{1, 2}", fmt.Sprintf("%#v", ints))

	_, err = arrutil.StringsToInts([]string{"a", "b"})
	is.Error(err)

	ints = arrutil.StringsAsInts([]string{"1", "2"})
	is.Equal("[]int{1, 2}", fmt.Sprintf("%#v", ints))
	is.Nil(arrutil.StringsAsInts([]string{"abc"}))
}

func TestStringsToSlice(t *testing.T) {
	is := assert.New(t)

	as := arrutil.StringsToSlice([]string{"1", "2"})
	is.Equal(`[]interface {}{"1", "2"}`, fmt.Sprintf("%#v", as))
}

func TestAnyToSlice(t *testing.T) {
	is := assert.New(t)

	sl, err := arrutil.AnyToSlice([]int{1, 2})
	is.NoError(err)
	is.Equal("[]interface {}{1, 2}", fmt.Sprintf("%#v", sl))

	_, err = arrutil.AnyToSlice(123)
	is.Error(err)
}

func TestToBoolSliceE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect []bool
		iserr  bool
	}{
		{[]bool{true, false, true}, []bool{true, false, true}, false},
		{[]interface{}{true, false, true}, []bool{true, false, true}, false},
		{[]int{1, 0, 1}, []bool{true, false, true}, false},
		{[]string{"true", "false", "true"}, []bool{true, false, true}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToBoolSliceE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToBoolSlice(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToIntSliceE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect []int
		iserr  bool
	}{
		{[]int{1, 3}, []int{1, 3}, false},
		{[]interface{}{1.2, 3.2}, []int{1, 3}, false},
		{[]string{"2", "3"}, []int{2, 3}, false},
		{[2]string{"2", "3"}, []int{2, 3}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToIntSliceE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToIntSlice(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToSliceE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect []interface{}
		iserr  bool
	}{
		{[]interface{}{1, 3}, []interface{}{1, 3}, false},
		{
			[]map[string]interface{}{{"k1": 1}, {"k2": 2}},
			[]interface{}{map[string]interface{}{"k1": 1}, map[string]interface{}{"k2": 2}},
			false,
		},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToSliceE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToSlice(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToStringSliceE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect []string
		iserr  bool
	}{
		{[]int{1, 2}, []string{"1", "2"}, false},
		{[]int8{int8(1), int8(2)}, []string{"1", "2"}, false},
		{[]int32{int32(1), int32(2)}, []string{"1", "2"}, false},
		{[]int64{int64(1), int64(2)}, []string{"1", "2"}, false},
		{[]float32{float32(1.01), float32(2.01)}, []string{"1.01", "2.01"}, false},
		{[]float64{float64(1.01), float64(2.01)}, []string{"1.01", "2.01"}, false},
		{[]string{"a", "b"}, []string{"a", "b"}, false},
		{[]interface{}{1, 3}, []string{"1", "3"}, false},
		{interface{}(1), []string{"1"}, false},
		{[]error{errors.New("a"), errors.New("b")}, []string{"a", "b"}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringSliceE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToStringSlice(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}
