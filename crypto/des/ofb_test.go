package des

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Decrypt data with key use DES OFB algorithm. Length of `key` param should be 8.
// func DesOfbDecrypt(encrypted, key []byte) []byte
func TestDesOfb(t *testing.T) {
	data := "hello"
	key := "abcdefgh"

	encrypted := cryptor.DesOfbEncrypt([]byte(data), []byte(key))
	decrypted := cryptor.DesOfbDecrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
