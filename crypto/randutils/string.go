package randutils

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

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
	var (
		result []byte
	)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var i int64 = 0
	for ; i < length; i++ {
		result = append(result, sources[r.Intn(len(sources))])
	}

	return string(result)
}

// RandomString 随机字符串
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
		return RandString2(64)
	}

	return fmt.Sprintf("%x", buf)
}
