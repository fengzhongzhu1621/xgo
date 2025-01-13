package aes

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Encrypt data with key use AES GCM algorithm.
// func AesGcmEncrypt(data, key []byte) []byte
// func AesGcmDecrypt(data, key []byte) []byte
func TestAesGcm(t *testing.T) {
	data := "hello"
	key := "abcdefghijklmnop"

	encrypted := cryptor.AesGcmEncrypt([]byte(data), []byte(key))
	decrypted := cryptor.AesGcmDecrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
