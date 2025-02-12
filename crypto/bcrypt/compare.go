package bcrypt

import "golang.org/x/crypto/bcrypt"

func CompareHashAndPassword(hashPassword, Password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(Password))

	return err == nil
}

func (pm *BCryptPasswordManager) CheckPassword(password, hash string) bool {
	return CompareHashAndPassword(password, hash)
}
