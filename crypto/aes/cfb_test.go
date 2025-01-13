package aes

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Encrypt data with key use AES CFB algorithm. Length of `key` param should be 16, 24 or 32.
// func AesCfbEncrypt(data, key []byte) []byte
// func AesCfbDecrypt(encrypted, key []byte) []byte
func TestAesCfb(t *testing.T) {
	data := "hello"
	key := "abcdefghijklmnop"

	encrypted := cryptor.AesCfbEncrypt([]byte(data), []byte(key))
	decrypted := cryptor.AesCfbDecrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
