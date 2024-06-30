package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

// GenerateID 获得随机字符串.
func GenerateID() string {
	base := 10
	return strconv.FormatInt(time.Now().UnixNano(), base)
}

// Md5 计算字符串的MD5值
// 同 echo -n "123456789" | md5sum.
func Md5(src string) string {
	md5ctx := md5.New()
	md5ctx.Write([]byte(src))
	cipher := md5ctx.Sum(nil)
	value := hex.EncodeToString(cipher)
	return value
}
