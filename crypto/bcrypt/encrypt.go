package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(encryptText string) (string, error) {
	hashStr, err := bcrypt.GenerateFromPassword([]byte(encryptText), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashStr), err
}

type BCryptPasswordManager struct {
	Cost int
}

func (pm *BCryptPasswordManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.Cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
