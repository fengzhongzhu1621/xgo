package aes

import (
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/stretchr/testify/assert"
)

func TestLancetAes(t *testing.T) {
	data := "hello"
	key := "abcdefghijklmnop"

	encrypted := cryptor.AesEcbEncrypt([]byte(data), []byte(key))
	decrypted := cryptor.AesEcbDecrypt(encrypted, []byte(key))

	assert.Equal(t, "hello", string(decrypted))

}
