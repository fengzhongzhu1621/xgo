package crypto

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/validator"
)

// func Utf8ToGbk(bs []byte) ([]byte, error)
func TestUtf8ToGbk(t *testing.T) {
	utf8Data := []byte("hello")
	gbkData, _ := convertor.Utf8ToGbk(utf8Data)

	fmt.Println(utf8.Valid(utf8Data))
	fmt.Println(validator.IsGBK(gbkData))

	// Output:
	// true
	// true
}

// func GbkToUtf8(bs []byte) ([]byte, error)
func TestGbkToUtf8(t *testing.T) {
	gbkData, _ := convertor.Utf8ToGbk([]byte("hello"))
	utf8Data, _ := convertor.GbkToUtf8(gbkData)

	fmt.Println(utf8.Valid(utf8Data))
	fmt.Println(string(utf8Data))

	// Output:
	// true
	// hello
}
