package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"

	"github.com/fengzhongzhu1621/xgo/str/bytesconv"
)

// reference: https://golang.org/src/crypto/cipher/example_test.go

const (
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256)

	ValidAES128KeySize int = 16
	ValidAES256KeySize int = 32

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.

	NonceByteSize int = 12
)

var (
	ErrInvalidKey   = errors.New("invalid key, should be 16 or 32 bytes")
	ErrInvalidNonce = errors.New("invalid nonce, should be 12 bytes")
)

type AESGcm struct {
	key   []byte
	nonce []byte
	// authenticated encryption with associated data (AEAD)
	aead cipher.AEAD
}

// NewAESGcm returns a new AES-GCM instance
func NewAESGcm(key []byte, nonce []byte) (aesGcm *AESGcm, err error) {
	// check key and nonce length
	if len(key) != ValidAES128KeySize && len(key) != ValidAES256KeySize {
		return nil, ErrInvalidKey
	}

	if len(nonce) != NonceByteSize {
		return nil, ErrInvalidNonce
	}

	// create AEAD
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	return &AESGcm{
		key:   key,
		nonce: nonce,
		aead:  aead,
	}, nil
}

// Encrypt encrypts plaintext
func (a *AESGcm) Encrypt(plaintext []byte) []byte {
	encryptedText := a.aead.Seal(plaintext[:0], a.nonce, plaintext, nil)
	return encryptedText
}

// Decrypt decrypts ciphertext
func (a *AESGcm) Decrypt(encryptedText []byte) ([]byte, error) {
	plaintext, err := a.aead.Open(nil, a.nonce, encryptedText, nil)
	return plaintext, err
}

// EncryptToString encrypts plaintext to string
func (a *AESGcm) EncryptToString(plaintext []byte) string {
	encryptedText := a.aead.Seal(plaintext[:0], a.nonce, plaintext, nil)
	return bytesconv.BytesToString(encryptedText)
}

// DecryptString decrypts ciphertext string
func (a *AESGcm) DecryptString(encryptedText string) ([]byte, error) {
	plaintext, err := a.aead.Open(nil, a.nonce, bytesconv.StringToBytes(encryptedText), nil)
	return plaintext, err
}
