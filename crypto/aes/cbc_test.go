package aes

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Encrypt data with key use AES CBC algorithm. Length of `key` param should be 16, 24 or 32.
// func AesCbcEncrypt(data, key []byte) []byte
// func AesCbcDecrypt(encrypted, key []byte) []byte
func TestAesCbc(t *testing.T) {
	data := "hello"
	key := "abcdefghijklmnop"

	encrypted := cryptor.AesCbcEncrypt([]byte(data), []byte(key))
	decrypted := cryptor.AesCbcDecrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
