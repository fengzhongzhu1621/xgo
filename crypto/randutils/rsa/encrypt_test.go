package rsa

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/cryptor"
)

// Create the rsa public and private key file in current directory.
// func GenerateRsaKey(keySize int, priKeyFile, pubKeyFile string) error
//
// Encrypt data with public key file useing ras algorithm.
// func RsaEncrypt(data []byte, pubKeyFileName string) []byte
//
// Decrypt data with private key file useing ras algorithm.
// func RsaDecrypt(data []byte, privateKeyFileName string) []byte
func TestRsaEncrypt(t *testing.T) {
	err := cryptor.GenerateRsaKey(4096, "rsa_private.pem", "rsa_public.pem")
	if err != nil {
		return
	}

	data := []byte("hello")
	encrypted := cryptor.RsaEncrypt(data, "rsa_public.pem")
	decrypted := cryptor.RsaDecrypt(encrypted, "rsa_private.pem")

	fmt.Println(string(decrypted))

	// Output:
	// hello
}

// Creates rsa private and public key.
// func GenerateRsaKeyPair(keySize int) (*rsa.PrivateKey, *rsa.PublicKey)
//
// Encrypts the given data with RSA-OAEP.
// func RsaEncryptOAEP(data []byte, label []byte, key rsa.PublicKey) ([]byte, error)
//
// Decrypts the data with RSA-OAEP.
// func RsaDecryptOAEP(ciphertext []byte, label []byte, key rsa.PrivateKey) ([]byte, error)
func TestRsaEncryptOAEP(t *testing.T) {
	// Creates rsa private and public key.
	pri, pub := cryptor.GenerateRsaKeyPair(1024)

	data := []byte("hello world")
	label := []byte("123456")

	encrypted, err := cryptor.RsaEncryptOAEP(data, label, *pub)
	if err != nil {
		return
	}

	decrypted, err := cryptor.RsaDecryptOAEP([]byte(encrypted), label, *pri)
	if err != nil {
		return
	}

	fmt.Println(string(decrypted))

	// Output:
	// hello world
}
