package rsa

import (
	"crypto"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Signs the data with RSA algorithm.
// func RsaSign(hash crypto.Hash, data []byte, privateKeyFileName string) ([]byte, error)
// func RsaVerifySign(hash crypto.Hash, data, signature []byte, pubKeyFileName string) error
func TestRsaSign(t *testing.T) {
	data := []byte("This is a test data for RSA signing")
	hash := crypto.SHA256

	privateKey := "./rsa_private_example.pem"
	publicKey := "./rsa_public_example.pem"

	signature, err := cryptor.RsaSign(hash, data, privateKey)
	if err != nil {
		return
	}

	err = cryptor.RsaVerifySign(hash, data, signature, publicKey)
	if err != nil {
		return
	}

}
