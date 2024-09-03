package randutils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v3"
	"github.com/oklog/ulid"
)

// NewUUID returns a new UUID Version 4.
func NewUUID() string {
	return uuid.New().String()
}

// NewShortUUID returns a new short UUID.
func NewShortUUID() string {
	return shortuuid.New()
}

// NewULID returns a new ULID.
func NewULID() string {
	return ulid.MustNew(ulid.Now(), rand.Reader).String()
}

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

// MD5Hash 计算字符串的MD5值
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
