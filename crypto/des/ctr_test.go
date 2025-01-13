package des

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Encrypt or decrypt data with key use DES CTR algorithm. Length of `key` param should be 8.
// func DesCtrCrypt(data, key []byte) []byte
// func DesCfbEncrypt(data, key []byte) []byte
func TestDesCtr(t *testing.T) {
	data := "hello"
	key := "abcdefgh"

	encrypted := cryptor.DesCtrCrypt([]byte(data), []byte(key))
	decrypted := cryptor.DesCtrCrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
