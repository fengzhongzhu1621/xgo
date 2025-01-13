package des

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Decrypt data with key use DES CBC algorithm. Length of `key` param should be 8.
// func DesCfbDecrypt(encrypted, key []byte) []byte
// func DesOfbEncrypt(data, key []byte) []byte
func TestDesCfb(t *testing.T) {
	data := "hello"
	key := "abcdefgh"

	encrypted := cryptor.DesCfbEncrypt([]byte(data), []byte(key))
	decrypted := cryptor.DesCfbDecrypt(encrypted, []byte(key))

	fmt.Println(string(decrypted))

	// Output:
	// hello
}
