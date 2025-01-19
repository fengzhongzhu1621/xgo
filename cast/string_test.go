package cast

import (
	"fmt"
	"html/template"
	"math/rand"
	"reflect"
	"testing"

	"github.com/gookit/goutil/byteutil"

	"github.com/gookit/goutil/arrutil"

	"github.com/duke-git/lancet/v2/convertor"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

var testString = "Albert Einstein: Logic will get you from A to B. Imagination will take you everywhere."
var testBytes = []byte(testString)

type foo struct {
	val string
}

func (x foo) String() string {
	return x.val
}

type fu struct {
	val string
}

func (x fu) Error() string {
	return x.val
}

// TestLancetStringToBytes 将字符串转换为字节切片而无需进行内存分配。
// func StringToBytes(str string) (b []byte)
func TestLancetStringToBytes(t *testing.T) {
	result1 := strutil.StringToBytes("abc")
	result2 := reflect.DeepEqual(result1, []byte{'a', 'b', 'c'})

	// [97 98 99]
	fmt.Println(result1)
	assert.Equal(t, true, result2)
}

// func ToChar(s string) []string
func TestToChar(t *testing.T) {
	result1 := convertor.ToChar("")
	result2 := convertor.ToChar("abc")
	result3 := convertor.ToChar("1 2#3")
	result4 := convertor.ToChar("你好")

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)

	// Output:
	// []
	// [a b c]
	// [1   2 # 3]
	// [你 好]
}

// func ToString(value any) string
func TestAnyToString(t *testing.T) {
	{
		result1 := convertor.ToString("")
		result2 := convertor.ToString(nil)
		result3 := convertor.ToString(0)
		result4 := convertor.ToString(1.23)
		result5 := convertor.ToString(true)
		result6 := convertor.ToString(false)
		result7 := convertor.ToString([]int{1, 2, 3})

		fmt.Println(result1)
		fmt.Println(result2)
		fmt.Println(result3)
		fmt.Println(result4)
		fmt.Println(result5)
		fmt.Println(result6)
		fmt.Println(result7)

		// Output:
		//
		//
		// 0
		// 1.23
		// true
		// false
		// [1,2,3]
	}

	{
		is := assert.New(t)
		arr := [2]string{"a", "b"}

		is.Equal("", arrutil.AnyToString(nil))
		is.Equal("[]", arrutil.AnyToString([]string{}))
		is.Equal("[a, b]", arrutil.AnyToString(arr))
		is.Equal("[a, b]", arrutil.AnyToString([]string{"a", "b"}))
		is.Equal("", arrutil.AnyToString("invalid"))
	}
}

// TestLancetBytesToString 将字节切片转换为字符串而无需进行内存分配。
// func BytesToString(bytes []byte) string
func TestBytesToString(t *testing.T) {
	{
		bytes := []byte{'a', 'b', 'c'}
		result := strutil.BytesToString(bytes)
		assert.Equal(t, "abc", result)
	}

	{
		data := make([]byte, 1024)
		for i := 0; i < 100; i++ {
			rand.Read(data)
			if rawBytesToStr(data) != BytesToString(data) {
				t.Fatal("don't match")
			}
		}

	}

	{
		assert.Equal(t, "123", byteutil.String([]byte("123")))
		assert.Equal(t, "123", byteutil.ToString([]byte("123")))
	}
}

func TestTruncateBytesToString(t *testing.T) {
	content := []byte("Hello, world!")
	truncatedStr := TruncateBytesToString(content, 5)
	assert.Equal(t, "Hello", string(truncatedStr))
}

func TestToStrings(t *testing.T) {
	is := assert.New(t)

	ss, err := arrutil.ToStrings([]int{1, 2})
	is.Nil(err)
	is.Equal(`[]string{"1", "2"}`, fmt.Sprintf("%#v", ss))

	ss = arrutil.MustToStrings([]int{1, 2})
	is.Equal(`[]string{"1", "2"}`, fmt.Sprintf("%#v", ss))

	ss = arrutil.MustToStrings([]any{1, 2})
	is.Equal(`[]string{"1", "2"}`, fmt.Sprintf("%#v", ss))

	ss = arrutil.SliceToStrings([]any{1, 2})
	is.Equal(`[]string{"1", "2"}`, fmt.Sprintf("%#v", ss))

	ss, err = arrutil.ToStrings("b")
	is.Nil(err)
	is.Equal(`[]string{"b"}`, fmt.Sprintf("%#v", ss))

	is.Empty(arrutil.AnyToStrings(234))
	is.Panics(func() {
		arrutil.MustToStrings(234)
	})

	_, err = arrutil.ToStrings([]any{[]int{1}, nil})
	is.Error(err)
}

func TestSliceToString(t *testing.T) {
	is := assert.New(t)
	is.Equal("[]", arrutil.SliceToString(nil))

	is.Equal("[]", arrutil.ToString[any](nil))
	is.Equal("[a,b]", arrutil.ToString([]any{"a", "b"}))

	is.Equal("[a,b]", arrutil.SliceToString("a", "b"))
}

func TestIntsToString(t *testing.T) {
	assert.Equal(t, "[]", arrutil.IntsToString([]int{}))
	assert.Equal(t, "[1,2,3]", arrutil.IntsToString([]int{1, 2, 3}))
}

// go test -v -run=none -bench=^BenchmarkBytesConv -benchmem=true

func BenchmarkBytesConvBytesToStrRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawBytesToStr(testBytes)
	}
}

func BenchmarkBytesConvBytesToStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BytesToString(testBytes)
	}
}

func BenchmarkBytesConvStrToBytesRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawStrToBytes(testString)
	}
}

func BenchmarkBytesConvStrToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringToBytes(testString)
	}
}

func TestToStringE(t *testing.T) {
	type Key struct {
		k string
	}
	key := &Key{"foo"}

	tests := []struct {
		input  interface{}
		expect string
		iserr  bool
	}{
		{int(8), "8", false},
		{int8(8), "8", false},
		{int16(8), "8", false},
		{int32(8), "8", false},
		{int64(8), "8", false},
		{uint(8), "8", false},
		{uint8(8), "8", false},
		{uint16(8), "8", false},
		{uint32(8), "8", false},
		{uint64(8), "8", false},
		{float32(8.31), "8.31", false},
		{float64(8.31), "8.31", false},
		{true, "true", false},
		{false, "false", false},
		{nil, "", false},
		{[]byte("one time"), "one time", false},
		{"one more time", "one more time", false},
		{template.HTML("one time"), "one time", false},
		{template.URL("http://somehost.foo"), "http://somehost.foo", false},
		{template.JS("(1+2)"), "(1+2)", false},
		{template.CSS("a"), "a", false},
		{template.HTMLAttr("a"), "a", false},
		// errors
		{testing.T{}, "", true},
		{key, "", true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToStringE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToString(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func TestStringerToString(t *testing.T) {
	var x foo
	x.val = "bar"
	assert.Equal(t, "bar", ToString(x))
}

func TestErrorToString(t *testing.T) {
	var x fu
	x.val = "bar"
	assert.Equal(t, "bar", ToString(x))
}
