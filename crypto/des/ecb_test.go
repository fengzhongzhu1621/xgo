package des

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Encrypt data with key use DES ECB algorithm. Length of `key` param should be 8.
// func DesEcbEncrypt(data, key []byte) []byte
// func DesEcbDecrypt(encrypted, key []byte) []byte
func TestDesEcb(t *testing.T) {
	data := "hello"
	key := "abcdefgh"

	encrypted := cryptor.DesEcbEncrypt([]byte(data), []byte(key))

	decrypted := cryptor.DesEcbDecrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
