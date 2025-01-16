package cast

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/convertor"
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

func TestTruncateBytes(t *testing.T) {
	content := []byte("Hello, world!")
	truncated := TruncateBytes(content, 5)
	assert.Equal(t, "Hello", string(truncated))
}
