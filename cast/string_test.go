package cast

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/duke-git/lancet/v2/convertor"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

func TestStringToBytes(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := RandStringBytesMaskImprSrcSB(64)
		if !bytes.Equal(rawStrToBytes(s), StringToBytes(s)) {
			t.Fatal("don't match")
		}
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
func TestToString(t *testing.T) {
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

var testString = "Albert Einstein: Logic will get you from A to B. Imagination will take you everywhere."
var testBytes = []byte(testString)

// go test -v

func TestBytesToString(t *testing.T) {
	data := make([]byte, 1024)
	for i := 0; i < 100; i++ {
		rand.Read(data)
		if rawBytesToStr(data) != BytesToString(data) {
			t.Fatal("don't match")
		}
	}
}

// TestLancetBytesToString 将字节切片转换为字符串而无需进行内存分配。
// func BytesToString(bytes []byte) string
func TestLancetBytesToString(t *testing.T) {
	bytes := []byte{'a', 'b', 'c'}
	result := strutil.BytesToString(bytes)
	assert.Equal(t, "abc", result)
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
