package cast

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gookit/goutil/arrutil"
	"github.com/stretchr/testify/assert"
)

// func ToFloat(value any) (float64, error)
func TestToFloat(t *testing.T) {
	result1, _ := convertor.ToFloat("")
	result2, err := convertor.ToFloat("abc")
	result3, _ := convertor.ToFloat("-1")
	result4, _ := convertor.ToFloat("-.11")
	result5, _ := convertor.ToFloat("1.23e3")
	result6, _ := convertor.ToFloat(true)

	fmt.Println(result1)
	fmt.Println(result2, err)
	fmt.Println(result3)
	fmt.Println(result4)
	fmt.Println(result5)
	fmt.Println(result6)

	// Output:
	// 0
	// 0 strconv.ParseFloat: parsing "": invalid syntax
	// -1
	// -0.11
	// 1230
	// 0
}

// func ToInt(value any) (int64, error)
func TestToInt(t *testing.T) {
	result1, _ := convertor.ToInt("123")
	result2, _ := convertor.ToInt("-123")
	result3, _ := convertor.ToInt(float64(12.3))
	result4, err := convertor.ToInt("abc")
	result5, _ := convertor.ToInt(true)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4, err)
	fmt.Println(result5)

	// Output:
	// 123
	// -123
	// 12
	// 0 strconv.ParseInt: parsing "": invalid syntax
	// 0
}

func TestToUintE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect uint
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8, false},
		{float64(8.31), 8, false},
		{true, 1, false},
		{false, 0, false},
		{"8", 8, false},
		{nil, 0, false},
		// errors
		{int(-8), 0, true},
		{int8(-8), 0, true},
		{int16(-8), 0, true},
		{int32(-8), 0, true},
		{int64(-8), 0, true},
		{float32(-8.31), 0, true},
		{float64(-8.31), 0, true},
		{"-8", 0, true},
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToUintE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test:
		v = ToUint(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToUint64E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect uint64
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8, false},
		{float64(8.31), 8, false},
		{true, 1, false},
		{false, 0, false},
		{"8", 8, false},
		{nil, 0, false},
		// errors
		{int(-8), 0, true},
		{int8(-8), 0, true},
		{int16(-8), 0, true},
		{int32(-8), 0, true},
		{int64(-8), 0, true},
		{float32(-8.31), 0, true},
		{float64(-8.31), 0, true},
		{"-8", 0, true},
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToUint64E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test:
		v = ToUint64(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToUint32E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect uint32
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8, false},
		{float64(8.31), 8, false},
		{true, 1, false},
		{false, 0, false},
		{"8", 8, false},
		{nil, 0, false},
		{int(-8), 0, true},
		{int8(-8), 0, true},
		{int16(-8), 0, true},
		{int32(-8), 0, true},
		{int64(-8), 0, true},
		{float32(-8.31), 0, true},
		{float64(-8.31), 0, true},
		{"-8", 0, true},
		// errors
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToUint32E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test:
		v = ToUint32(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToUint16E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect uint16
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8, false},
		{float64(8.31), 8, false},
		{true, 1, false},
		{false, 0, false},
		{"8", 8, false},
		{nil, 0, false},
		// errors
		{int(-8), 0, true},
		{int8(-8), 0, true},
		{int16(-8), 0, true},
		{int32(-8), 0, true},
		{int64(-8), 0, true},
		{float32(-8.31), 0, true},
		{float64(-8.31), 0, true},
		{"-8", 0, true},
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToUint16E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToUint16(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToUint8E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect uint8
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8, false},
		{float64(8.31), 8, false},
		{true, 1, false},
		{false, 0, false},
		{"8", 8, false},
		{nil, 0, false},
		// errors
		{int(-8), 0, true},
		{int8(-8), 0, true},
		{int16(-8), 0, true},
		{int32(-8), 0, true},
		{int64(-8), 0, true},
		{float32(-8.31), 0, true},
		{float64(-8.31), 0, true},
		{"-8", 0, true},
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToUint8E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToUint8(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToIntE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect int
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8, false},
		{float64(8.31), 8, false},
		{true, 1, false},
		{false, 0, false},
		{"8", 8, false},
		{nil, 0, false},
		// errors
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToIntE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToInt(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToInt64E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect int64
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8, false},
		{float64(8.31), 8, false},
		{true, 1, false},
		{false, 0, false},
		{"8", 8, false},
		{nil, 0, false},
		// errors
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToInt64E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToInt64(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToInt32E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect int32
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8, false},
		{float64(8.31), 8, false},
		{true, 1, false},
		{false, 0, false},
		{"8", 8, false},
		{nil, 0, false},
		// errors
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToInt32E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToInt32(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToInt16E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect int16
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8, false},
		{float64(8.31), 8, false},
		{true, 1, false},
		{false, 0, false},
		{"8", 8, false},
		{nil, 0, false},
		// errors
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToInt16E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToInt16(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToInt8E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect int8
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8, false},
		{float64(8.31), 8, false},
		{true, 1, false},
		{false, 0, false},
		{"8", 8, false},
		{nil, 0, false},
		// errors
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToInt8E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToInt8(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToFloat64E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect float64
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8), 8, false},
		{float64(8.31), 8.31, false},
		{"8", 8, false},
		{true, 1, false},
		{false, 0, false},
		// errors
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToFloat64E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToFloat64(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestToFloat32E(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect float32
		iserr  bool
	}{
		{int(8), 8, false},
		{int8(8), 8, false},
		{int16(8), 8, false},
		{int32(8), 8, false},
		{int64(8), 8, false},
		{uint(8), 8, false},
		{uint8(8), 8, false},
		{uint16(8), 8, false},
		{uint32(8), 8, false},
		{uint64(8), 8, false},
		{float32(8.31), 8.31, false},
		{float64(8.31), 8.31, false},
		{"8", 8, false},
		{true, 1, false},
		{false, 0, false},
		// errors
		{"test", 0, true},
		{testing.T{}, 0, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToFloat32E(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToFloat32(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

// TestToBigInt Converts an integer of any supported type (int, int64, uint64, etc.) to *big.Int
// func ToBigInt[T any](v T) (*big.Int, error)
func TestToBigInt(t *testing.T) {
	n := 9876543210
	bigInt, _ := convertor.ToBigInt(n)

	fmt.Println(bigInt)
	// Output:
	// 9876543210
}

func TestToInt64s(t *testing.T) {
	is := assert.New(t)

	ints, err := arrutil.ToInt64s([]string{"1", "2"})
	is.Nil(err)
	is.Equal("[]int64{1, 2}", fmt.Sprintf("%#v", ints))

	ints = arrutil.MustToInt64s([]string{"1", "2"})
	is.Equal("[]int64{1, 2}", fmt.Sprintf("%#v", ints))

	ints = arrutil.MustToInt64s([]any{"1", "2"})
	is.Equal("[]int64{1, 2}", fmt.Sprintf("%#v", ints))

	ints = arrutil.SliceToInt64s([]any{"1", "2"})
	is.Equal("[]int64{1, 2}", fmt.Sprintf("%#v", ints))

	_, err = arrutil.ToInt64s([]string{"a", "b"})
	is.Error(err)
}

func TestIndirectPointers(t *testing.T) {
	x := 13
	y := &x
	z := &y

	assert.Equal(t, ToInt(y), 13)
	assert.Equal(t, ToInt(z), 13)
}

func TestGconv(t *testing.T) {
	i := 123
	fmt.Printf("%10s %v\n", "Int:", gconv.Int(i))
	fmt.Printf("%10s %v\n", "Int8:", gconv.Int8(i))
	fmt.Printf("%10s %v\n", "Int16:", gconv.Int16(i))
	fmt.Printf("%10s %v\n", "Int32:", gconv.Int32(i))
	fmt.Printf("%10s %v\n", "Int64:", gconv.Int64(i))
	fmt.Printf("%10s %v\n", "Uint:", gconv.Uint(i))
	fmt.Printf("%10s %v\n", "Float32:", gconv.Float32(i))
	fmt.Printf("%10s %v\n", "Float64:", gconv.Float64(i))
	fmt.Printf("%10s %v\n", "Bool:", gconv.Bool(i))
	fmt.Printf("%10s %v\n", "String:", gconv.String(i))

	//      Int: 123
	//     Int8: 123
	//    Int16: 123
	//    Int32: 123
	//    Int64: 123
	//     Uint: 123
	//  Float32: 123
	//  Float64: 123
	//     Bool: true
	//   String: 123
}
