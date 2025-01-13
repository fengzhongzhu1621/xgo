package aes

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Encrypt data with key use AES OFB algorithm. Length of `key` param should be 16, 24 or 32.
// func AesOfbEncrypt(data, key []byte) []byte
// func AesOfbDecrypt(encrypted, key []byte) []byte
func TestAesOfb(t *testing.T) {
	data := "hello"
	key := "abcdefghijklmnop"

	encrypted := cryptor.AesOfbEncrypt([]byte(data), []byte(key))
	decrypted := cryptor.AesCfbDecrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
