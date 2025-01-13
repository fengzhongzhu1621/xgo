package crypto

import (
	"encoding/base64"
	"errors"
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"

	"github.com/duke-git/lancet/v2/convertor"

	"github.com/stretchr/testify/assert"
)

func StdBase64(s string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(s))
}

func TestBase64(t *testing.T) {
	msg := []byte("Hello world！ URL safe 编码，相当于替换掉字符串中的特殊字符，+ 和 /")

	// 1、标准编码
	encoded := base64.StdEncoding.EncodeToString(msg)
	assert.Equal(t, "SGVsbG8gd29ybGTvvIEgVVJMIHNhZmUg57yW56CB77yM55u45b2T5LqO5pu/5o2i5o6J5a2X56ym5Liy5Lit55qE54m55q6K5a2X56ym77yMKyDlkowgLw==", encoded)
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	assert.Equal(t, msg, decoded)

	// 2、常规编码，末尾不补 =
	encoded = base64.RawStdEncoding.EncodeToString(msg)
	assert.Equal(t, "SGVsbG8gd29ybGTvvIEgVVJMIHNhZmUg57yW56CB77yM55u45b2T5LqO5pu/5o2i5o6J5a2X56ym5Liy5Lit55qE54m55q6K5a2X56ym77yMKyDlkowgLw", encoded)
	decoded, _ = base64.RawStdEncoding.DecodeString(encoded)
	assert.Equal(t, msg, decoded)

	// 3、URL safe 编码	, 替换掉字符串中的特殊字符，+ 和 /
	encoded = base64.URLEncoding.EncodeToString(msg)
	assert.Equal(t, "SGVsbG8gd29ybGTvvIEgVVJMIHNhZmUg57yW56CB77yM55u45b2T5LqO5pu_5o2i5o6J5a2X56ym5Liy5Lit55qE54m55q6K5a2X56ym77yMKyDlkowgLw==", encoded)
	decoded, _ = base64.URLEncoding.DecodeString(encoded)
	assert.Equal(t, msg, decoded)

	// 4、URL safe 编码, 替换掉字符串中的特殊字符，+ 和 /，末尾不补 =
	encoded = base64.RawURLEncoding.EncodeToString(msg)
	assert.Equal(t, "SGVsbG8gd29ybGTvvIEgVVJMIHNhZmUg57yW56CB77yM55u45b2T5LqO5pu_5o2i5o6J5a2X56ym5Liy5Lit55qE54m55q6K5a2X56ym77yMKyDlkowgLw", encoded)
	decoded, _ = base64.RawURLEncoding.DecodeString(encoded)
	assert.Equal(t, msg, decoded)
}

// Encode string with base64 encoding.
// func Base64StdEncode(s string) string
// func Base64StdDecode(s string) string
func TestBase64StdEncode(t *testing.T) {
	base64Str := cryptor.Base64StdEncode("hello")
	fmt.Println(base64Str)

	// Output:
	// aGVsbG8=

	str := cryptor.Base64StdDecode("aGVsbG8=")
	fmt.Println(str)

	// Output:
	// hello
}

// func ToStdBase64(value any) string
func TestToStdBase64(t *testing.T) {
	afterEncode := convertor.ToStdBase64(nil)
	fmt.Println(afterEncode)

	afterEncode = convertor.ToStdBase64("")
	fmt.Println(afterEncode)

	stringVal := "hello"
	afterEncode = convertor.ToStdBase64(stringVal)
	fmt.Println(afterEncode)

	byteSliceVal := []byte("hello")
	afterEncode = convertor.ToStdBase64(byteSliceVal)
	fmt.Println(afterEncode)

	intVal := 123
	afterEncode = convertor.ToStdBase64(intVal)
	fmt.Println(afterEncode)

	mapVal := map[string]any{"a": "hi", "b": 2, "c": struct {
		A string
		B int
	}{"hello", 3}}
	afterEncode = convertor.ToStdBase64(mapVal)
	fmt.Println(afterEncode)

	floatVal := 123.456
	afterEncode = convertor.ToStdBase64(floatVal)
	fmt.Println(afterEncode)

	boolVal := true
	afterEncode = convertor.ToStdBase64(boolVal)
	fmt.Println(afterEncode)

	errVal := errors.New("err")
	afterEncode = convertor.ToStdBase64(errVal)
	fmt.Println(afterEncode)

	// Output:
	//
	//
	// aGVsbG8=
	// aGVsbG8=
	// MTIz
	// eyJhIjoiaGkiLCJiIjoyLCJjIjp7IkEiOiJoZWxsbyIsIkIiOjN9fQ==
	// MTIzLjQ1Ng==
	// dHJ1ZQ==
	// ZXJy
}

// func ToUrlBase64(value any) string
func TestToUrlBase64(t *testing.T) {
	afterEncode := convertor.ToUrlBase64(nil)
	fmt.Println(afterEncode)

	stringVal := "hello"
	afterEncode = convertor.ToUrlBase64(stringVal)
	fmt.Println(afterEncode)

	byteSliceVal := []byte("hello")
	afterEncode = convertor.ToUrlBase64(byteSliceVal)
	fmt.Println(afterEncode)

	intVal := 123
	afterEncode = convertor.ToUrlBase64(intVal)
	fmt.Println(afterEncode)

	mapVal := map[string]any{"a": "hi", "b": 2, "c": struct {
		A string
		B int
	}{"hello", 3}}
	afterEncode = convertor.ToUrlBase64(mapVal)
	fmt.Println(afterEncode)

	floatVal := 123.456
	afterEncode = convertor.ToUrlBase64(floatVal)
	fmt.Println(afterEncode)

	boolVal := true
	afterEncode = convertor.ToUrlBase64(boolVal)
	fmt.Println(afterEncode)

	errVal := errors.New("err")
	afterEncode = convertor.ToUrlBase64(errVal)
	fmt.Println(afterEncode)

	// Output:
	//
	// aGVsbG8=
	// aGVsbG8=
	// MTIz
	// eyJhIjoiaGkiLCJiIjoyLCJjIjp7IkEiOiJoZWxsbyIsIkIiOjN9fQ==
	// MTIzLjQ1Ng==
	// dHJ1ZQ==
	// ZXJy
}

// func ToRawStdBase64(value any) string
func TestToRawStdBase64(t *testing.T) {

	stringVal := "hello"
	afterEncode := convertor.ToRawStdBase64(stringVal)
	fmt.Println(afterEncode)

	byteSliceVal := []byte("hello")
	afterEncode = convertor.ToRawStdBase64(byteSliceVal)
	fmt.Println(afterEncode)

	intVal := 123
	afterEncode = convertor.ToRawStdBase64(intVal)
	fmt.Println(afterEncode)

	mapVal := map[string]any{"a": "hi", "b": 2, "c": struct {
		A string
		B int
	}{"hello", 3}}
	afterEncode = convertor.ToRawStdBase64(mapVal)
	fmt.Println(afterEncode)

	floatVal := 123.456
	afterEncode = convertor.ToRawStdBase64(floatVal)
	fmt.Println(afterEncode)

	boolVal := true
	afterEncode = convertor.ToRawStdBase64(boolVal)
	fmt.Println(afterEncode)

	errVal := errors.New("err")
	afterEncode = convertor.ToRawStdBase64(errVal)
	fmt.Println(afterEncode)

	// Output:
	// aGVsbG8
	// aGVsbG8
	// MTIz
	// eyJhIjoiaGkiLCJiIjoyLCJjIjp7IkEiOiJoZWxsbyIsIkIiOjN9fQ
	// MTIzLjQ1Ng
	// dHJ1ZQ
	// ZXJy
}

// func ToRawUrlBase64(value any) string
func TestToRawUrlBase64(t *testing.T) {

	stringVal := "hello"
	afterEncode := convertor.ToRawUrlBase64(stringVal)
	fmt.Println(afterEncode)

	byteSliceVal := []byte("hello")
	afterEncode = convertor.ToRawUrlBase64(byteSliceVal)
	fmt.Println(afterEncode)

	intVal := 123
	afterEncode = convertor.ToRawUrlBase64(intVal)
	fmt.Println(afterEncode)

	mapVal := map[string]any{"a": "hi", "b": 2, "c": struct {
		A string
		B int
	}{"hello", 3}}
	afterEncode = convertor.ToRawUrlBase64(mapVal)
	fmt.Println(afterEncode)

	floatVal := 123.456
	afterEncode = convertor.ToRawUrlBase64(floatVal)
	fmt.Println(afterEncode)

	boolVal := true
	afterEncode = convertor.ToRawUrlBase64(boolVal)
	fmt.Println(afterEncode)

	errVal := errors.New("err")
	afterEncode = convertor.ToRawUrlBase64(errVal)
	fmt.Println(afterEncode)

	// Output:
	// aGVsbG8
	// aGVsbG8
	// MTIz
	// eyJhIjoiaGkiLCJiIjoyLCJjIjp7IkEiOiJoZWxsbyIsIkIiOjN9fQ
	// MTIzLjQ1Ng
	// dHJ1ZQ
	// ZXJy
}
