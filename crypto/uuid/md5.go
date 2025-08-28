package uuid

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5 计算字符串的MD5值
// 同 echo -n "123456789" | md5sum.
func Md5(src string) string {
	md5ctx := md5.New()
	md5ctx.Write([]byte(src))
	cipher := md5ctx.Sum(nil)
	value := hex.EncodeToString(cipher)
	return value
}

// MD5Hash 计算字符串的MD5值
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
