package randutils

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"math/rand"
	"sync/atomic"
	"time"
)

var nonceCounter uint64

func generateRandom(size int) []byte {
	nonce := atomic.AddUint64(&nonceCounter, 1)
	b := make([]byte, 13+8+13)
	rand.Read(b[:13])
	copy(b[13:21], u64ToBytes(nonce))
	rand.Read(b[21:])
	hash := sha256.Sum256(b)
	return hash[:size]
}

func u64ToBytes(n uint64) []byte {
	out := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		out[i] = byte(n)
		n >>= 8
	}
	return out
}

// RandomString 生成一个指定长度的随机字符串
func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// RandString2 生成一个指定长度的随机字符串
func RandString2(length int64) string {
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var i int64 = 0
	for ; i < length; i++ {
		result = append(result, sources[r.Intn(len(sources))])
	}

	return string(result)
}

// RandomString3 随机字符串
func RandomString3(size int) string {
	buf := make([]byte, size, size)
	max := big.NewInt(int64(chLen))
	for i := 0; i < size; i++ {
		random, err := crand.Int(crand.Reader, max)
		if err != nil {
			rand.Seed(time.Now().UnixNano())
			buf[i] = character[rand.Intn(chLen)]
			continue
		}
		buf[i] = character[random.Int64()]
	}

	return string(buf)
}

// RandAuthToken 随机生成 token
func RandAuthToken() string {
	buf := make([]byte, 32)
	_, err := crand.Read(buf)
	if err != nil {
		return RandString2(32)
	}

	return fmt.Sprintf("%x", buf)
}

func GenerateSecureID() string {
	b := make([]byte, 32)
	if _, err := crand.Read(b); err != nil {
		panic(err) // 或者返回错误，视具体需求而定
	}
	return base64.URLEncoding.EncodeToString(b)
}

func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := crand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateRandomPassword(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	bytes := make([]byte, length)

	if _, err := crand.Read(bytes); err != nil {
		return "", err
	}
	count := len(chars)

	// 填充字符
	for i, b := range bytes {
		bytes[i] = chars[b%byte(count)]
	}

	return string(bytes), nil
}
