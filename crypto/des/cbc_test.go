package des

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Encrypt data with key use DES CBC algorithm. Length of `key` param should be 8.
// func DesCbcEncrypt(data, key []byte) []byte
// func DesCbcDecrypt(encrypted, key []byte) []byte
func TestDesCbc(t *testing.T) {
	data := "hello"
	key := "abcdefgh"

	encrypted := cryptor.DesCbcEncrypt([]byte(data), []byte(key))
	decrypted := cryptor.DesCbcDecrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
