package cast

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gookit/goutil/byteutil"
	"github.com/stretchr/testify/assert"
)

// func EncodeByte(data any) ([]byte, error)
func TestEncodeByte(t *testing.T) {
	byteData, _ := convertor.EncodeByte("abc")
	fmt.Println(byteData)

	// Output:
	// [6 12 0 3 97 98 99]
}

// func DecodeByte(data []byte, target any) error
func TestDecodeByte(t *testing.T) {
	var result string
	byteData := []byte{6, 12, 0, 3, 97, 98, 99}

	err := convertor.DecodeByte(byteData, &result)
	if err != nil {
		return
	}

	fmt.Println(result)

	// Output:
	// abc
}

func TestStringToBytes(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := RandStringBytesMaskImprSrcSB(64)
		if !bytes.Equal(rawStrToBytes(s), StringToBytes(s)) {
			t.Fatal("don't match")
		}
	}
}

// func ToBytes(data any) ([]byte, error)
func TestToBytes(t *testing.T) {
	{
		bytesData, err := convertor.ToBytes("abc")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(bytesData)

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
