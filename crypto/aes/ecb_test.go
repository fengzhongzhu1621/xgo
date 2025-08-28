package aes

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Encrypt data with key use AES ECB algorithm. Length of `key` param should be 16, 24 or 32.
// func AesEcbEncrypt(data, key []byte) []byte
// func AesEcbDecrypt(encrypted, key []byte) []byte
func TestAesEcb(t *testing.T) {
	data := "hello"
	key := "abcdefghijklmnop"

	encrypted := cryptor.AesEcbEncrypt([]byte(data), []byte(key))
	decrypted := cryptor.AesEcbDecrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
