package crypto

type Crypto interface {
	Encrypt(plaintext []byte) []byte
	Decrypt(encryptedText []byte) ([]byte, error)

	EncryptToString(plaintext []byte) string
	DecryptString(encryptedText string) ([]byte, error)
}
