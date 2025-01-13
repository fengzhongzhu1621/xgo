package aes

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Encrypt or decrypt data with key use AES CTR algorithm. Length of `key` param should be 16, 24 or 32.
// func AesCtrCrypt(data, key []byte) []byte
// func AesCfbEncrypt(data, key []byte) []byte
func TestAesCtr(t *testing.T) {
	data := "hello"
	key := "abcdefghijklmnop"

	encrypted := cryptor.AesCtrCrypt([]byte(data), []byte(key))
	decrypted := cryptor.AesCtrCrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
