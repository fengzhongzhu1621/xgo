package randutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"

	"github.com/duke-git/lancet/v2/fileutil"
)

// returns file sha value, param `shaType` should be 1, 256 or 512.
// func Sha(filepath string, shaType ...int) (string, error)
func TestSha(t *testing.T) {
	fname := "./test.txt"
	fileutil.CreateFile(fname)

	f, _ := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0777)
	defer f.Close()

	f.WriteString("hello\nworld")

	sha1, err := fileutil.Sha(fname, 1)
	sha256, _ := fileutil.Sha(fname, 256)
	sha512, _ := fileutil.Sha(fname, 512)

	fmt.Println(sha1)
	fmt.Println(sha256)
	fmt.Println(sha512)
	fmt.Println(err)

	os.Remove(fname)
}

// func Sha1(str string) string
func TestSha1(t *testing.T) {
	str := "hello"

	result := cryptor.Sha1(str)
	fmt.Println(result)

	// Output:
	// aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
}

// func Sha1WithBase64(str string) string
func TestSha1WithBase64(t *testing.T) {
	result := cryptor.Sha1WithBase64("hello")
	fmt.Println(result)

	// Output:
	// qvTGHdzF6KLavt4PO0gs2a6pQ00=
}

// Get the sha1 hmac hash of string.
// func HmacSha1(str, key string) string
func TestHmacSha1(t *testing.T) {
	str := "hello"
	key := "12345"

	hms := cryptor.HmacSha1(str, key)
	fmt.Println(hms)

	// Output:
	// 5c6a9db0cccb92e36ed0323fd09b7f936de9ace0
}

// func Sha256(str string) string
func TestSha256(t *testing.T) {
	str := "hello"

	result := cryptor.Sha256(str)
	fmt.Println(result)

	// Output:
	// 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
}

// func Sha256WithBase64(str string) string
func TestSha256WithBase64(t *testing.T) {
	result := cryptor.Sha256WithBase64("hello")
	fmt.Println(result)

	// Output:
	// LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ=
}

// func Sha512(str string) string
func TestSha512(t *testing.T) {
	str := "hello"

	result := cryptor.Sha512(str)
	fmt.Println(result)

	// Output:
	// 9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043
}

// func Sha512WithBase64(str string) string
func TestSha512WithBase64(t *testing.T) {
	result := cryptor.Sha512WithBase64("hello")
	fmt.Println(result)

	// Output:
	// m3HSJL1i83hdltRq0+o9czGb+8KJDKra4t/3JRlnPKcjI8PZm6XBHXx6zG4UuMXaDEZjR1wuXDre9G9zvN7AQw==
}

// func HmacSha1WithBase64(str, key string) string
func TestHmacSha1WithBase64(t *testing.T) {
	str := "hello"
	key := "12345"

	hms := cryptor.HmacSha1WithBase64(str, key)
	fmt.Println(hms)

	// Output:
	// XGqdsMzLkuNu0DI/0Jt/k23prOA=
}

// func HmacSha256(str, key string) string
func TestsHmacSha256(t *testing.T) {
	str := "hello"
	key := "12345"

	hms := cryptor.HmacSha256(str, key)
	fmt.Println(hms)

	// Output:
	// 315bb93c4e989862ba09cb62e05d73a5f376cb36f0d786edab0c320d059fde75
}

// func HmacSha256WithBase64(str, key string) string
func TestHmacSha256WithBase64(t *testing.T) {
	str := "hello"
	key := "12345"

	hms := cryptor.HmacSha256WithBase64(str, key)
	fmt.Println(hms)

	// Output:
	// MVu5PE6YmGK6Ccti4F1zpfN2yzbw14btqwwyDQWf3nU=
}

// func HmacSha512(str, key string) string
func TestHmacSha512(t *testing.T) {
	str := "hello"
	key := "12345"

	hms := cryptor.HmacSha512(str, key)
	fmt.Println(hms)

	// Output:
	// dd8f1290a9dd23d354e2526d9a2e9ce8cffffdd37cb320800d1c6c13d2efc363288376a196c5458daf53f8e1aa6b45a6d856303d5c0a2064bff9785861d48cfc
}

// func HmacSha512WithBase64(str, key string) string
func TestHmacSha512WithBase64(t *testing.T) {
	str := "hello"
	key := "12345"

	hms := cryptor.HmacSha512WithBase64(str, key)
	fmt.Println(hms)

	// Output:
	// 3Y8SkKndI9NU4lJtmi6c6M///dN8syCADRxsE9Lvw2Mog3ahlsVFja9T+OGqa0Wm2FYwPVwKIGS/+XhYYdSM/A==
}
