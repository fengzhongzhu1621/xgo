package cast

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"reflect"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gookit/goutil/byteutil"
	"github.com/stretchr/testify/assert"
)

func TestStringToBytes(t *testing.T) {
	{
		// func EncodeByte(data any) ([]byte, error)
		byteData, _ := convertor.EncodeByte("abc")
		fmt.Println(string(byteData))

		// Output:
		// [6 12 0 3 97 98 99]
	}

	{
		// TestLancetStringToBytes 将字符串转换为字节切片而无需进行内存分配。
		// func StringToBytes(str string) (b []byte)
		result1 := strutil.StringToBytes("abc")
		result2 := reflect.DeepEqual(result1, []byte{'a', 'b', 'c'})

		// [97 98 99]
		fmt.Println(string(result1))
		assert.Equal(t, true, result2)
	}

	{
		for i := 0; i < 100; i++ {
			s := RandStringBytesMaskImprSrcSB(64)
			if !bytes.Equal(rawStrToBytes(s), StringToBytes(s)) {
				t.Fatal("don't match")
			}
		}
	}

	{
		// func ToBytes(data any) ([]byte, error)
		bytesData, err := convertor.ToBytes("abc")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(bytesData))

		// Output:
		// [97 98 99]
	}

	{
		tests := []struct {
			v   any
			exp []byte
			ok  bool
		}{
			{nil, nil, true},
			{123, []byte("123"), true},
			{int8(123), []byte("123"), true},
			{int16(123), []byte("123"), true},
			{int32(123), []byte("123"), true},
			{int64(123), []byte("123"), true},
			{uint(123), []byte("123"), true},
			{uint8(123), []byte("123"), true},
			{uint16(123), []byte("123"), true},
			{uint32(123), []byte("123"), true},
			{uint64(123), []byte("123"), true},
			{float32(123), []byte("123"), true},
			{float64(123), []byte("123"), true},
			{[]byte("123"), []byte("123"), true},
			{"123", []byte("123"), true},
			{true, []byte("true"), true},
			// special
			{time.Duration(123), []byte("123"), true},
			{fs.ModePerm, []byte("-rwxrwxrwx"), true},
			{errors.New("error msg"), []byte("error msg"), true},
			// failed
			{[]string{"123"}, nil, false},
		}

		for _, item := range tests {
			bs, err := byteutil.ToBytes(item.v)
			assert.Equal(t, item.ok, err == nil)
			assert.Equal(t, item.exp, bs, "real value: %v", item.v)
		}

		// SafeBytes
		bs := byteutil.SafeBytes([]string{"123"})
		assert.Equal(t, []byte("[123]"), bs)
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
