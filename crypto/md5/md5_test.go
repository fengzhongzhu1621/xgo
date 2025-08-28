package md5

import (
	"errors"
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/gookit/goutil/byteutil"
	"github.com/stretchr/testify/assert"
)

func TestGetMD5Hash(t *testing.T) {
	text := "xxx"
	expected := "f561aaf6ef0bf14d4208bb46a4ccb3ad"
	result := GetMD5Hash(text)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetBytesMD5Hash(t *testing.T) {
	text := []byte("xxx")
	expected := "f561aaf6ef0bf14d4208bb46a4ccb3ad"
	result := GetBytesMD5Hash(text)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetSliceMD5Hash(t *testing.T) {
	strSlice := []string{"xxx", "yyy"}
	expectedStr := "5f221cf63a70ca156f0fe1058e7f250b"
	resultStr, err := GetSliceMD5Hash(strSlice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resultStr != expectedStr {
		t.Errorf("Expected %v, but got %v", expectedStr, resultStr)
	}
	intSlice := []int64{1, 2, 3}
	expectedInt := "55b84a9d317184fe61224bfb4a060fb0"
	resultInt, err := GetSliceMD5Hash(intSlice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resultInt != expectedInt {
		t.Errorf("Expected %v, but got %v", expectedInt, resultInt)
	}
	invalidSlice := []float64{1.1, 2.2, 3.3}
	_, err = GetSliceMD5Hash(invalidSlice)
	assert.Error(t, err)
	assert.Equal(t, errors.New("illegal type").Error(), err.Error())
}

// func HmacMd5(str, key string) string
func TestHmacMd5(t *testing.T) {
	str := "hello"
	key := "12345"

	hms := cryptor.HmacMd5(str, key)
	fmt.Println(hms)

	// Output:
	// e834306eab892d872525d4918a7a639a
}

// Get the md5 hmac hash of base64 string.
// func HmacMd5WithBase64(str, key string) string
func TestHmacMd5WithBase64(t *testing.T) {
	str := "hello"
	key := "12345"

	hms := cryptor.HmacMd5WithBase64(str, key)
	fmt.Println(hms)

	// Output:
	// 6DQwbquJLYclJdSRinpjmg==
}

// func Md5String(s string) string
func TestMd5String(t *testing.T) {
	str := "hello"

	md5Str := cryptor.Md5String(str)
	fmt.Println(md5Str)

	// Output:
	// 5d41402abc4b2a76b9719d911017c592
}

// func Md5StringWithBase64(s string) string
func TestMd5StringWithBase64(t *testing.T) {
	md5Str := cryptor.Md5StringWithBase64("hello")
	fmt.Println(md5Str)

	// Output:
	// XUFAKrxLKna5cZ2REBfFkg==
}

// func Md5Byte(data []byte) string
func TestMd5Byte(t *testing.T) {
	md5Str := cryptor.Md5Byte([]byte{'a'})
	fmt.Println(md5Str)

	// Output:
	// 0cc175b9c0f1b6a831c399e269772661
}

// func Md5ByteWithBase64(data []byte) string
func TestMd5ByteWithBase64(t *testing.T) {
	md5Str := cryptor.Md5ByteWithBase64([]byte("hello"))
	fmt.Println(md5Str)

	// Output:
	// XUFAKrxLKna5cZ2REBfFkg==
}

// func Md5File(filepath string) (string, error)
func TestMd5File(t *testing.T) {
	s, _ := cryptor.Md5File("./main.go")
	fmt.Println(s)
}

func TestMd5(t *testing.T) {
	assert.NotEmpty(t, byteutil.Md5("abc"))
	assert.NotEmpty(t, byteutil.Md5([]int{12, 34}))

	assert.Equal(t, "202cb962ac59075b964b07152d234b70", string(byteutil.Md5("123")))
	assert.Equal(t, "900150983cd24fb0d6963f7d28e17f72", string(byteutil.Md5("abc")))

	// short md5
	assert.Equal(t, "ac59075b964b0715", string(byteutil.ShortMd5("123")))
	assert.Equal(t, "3cd24fb0d6963f7d", string(byteutil.ShortMd5("abc")))
}
