package randutils

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	sources     = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

var (
	character = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	chLen     = len(character)
)

var tmpMu sync.Mutex

// Int returns a non-negative pseudo-random int.
func Int() int { return pseudo.Int() }

// Intn returns, as an int, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func Intn(n int) int { return pseudo.Intn(n) }

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func Int63n(n int64) int64 { return pseudo.Int63n(n) }

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n).
func Perm(n int) []int { return pseudo.Perm(n) }

// Seed uses the provided seed value to initialize the default Source to a
// deterministic state. If Seed is not called, the generator behaves as if
// seeded by Seed(1).
func Seed(n int64) { pseudo.Seed(n) }

var pseudo = rand.New(&source{src: rand.NewSource(1)})

type source struct {
	src rand.Source
	mu  sync.Mutex
}

// Int63 支持多线程安全.
func (s *source) Int63() int64 {
	s.mu.Lock()
	n := s.src.Int63()
	s.mu.Unlock()
	return n
}

// Seed 支持多线程安全.
func (s *source) Seed(seed int64) {
	s.mu.Lock()
	s.src.Seed(seed)
	s.mu.Unlock()
}

// Shuffle pseudo-randomizes the order of elements.
// n is the number of elements.
// swap swaps the elements with indexes i and j.
func Shuffle(n int, swap func(i, j int)) { pseudo.Shuffle(n, swap) }

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

// RandAuthToken 随机生成 token
func RandAuthToken() string {
	buf := make([]byte, 32)
	_, err := crand.Read(buf)
	if err != nil {
		return RandString2(64)
	}

	return fmt.Sprintf("%x", buf)
}

// Random 随机数据
func Random(size int) []byte {
	buf := make([]byte, size)

	if _, err := crand.Read(buf); err != nil {
		rand.Seed(time.Now().UnixNano())
		rand.Read(buf)
	}
	return buf
}

// RandomInt 生成一个指定范围的随机整数
func RandomInt(min, max int) int {
	// crand.Int 是一个生成随机数的函数，它使用加密安全的随机数生成器（CSPRNG）
	// crand.Reader 是一个全局、共享的加密安全随机数生成器
	// big.NewInt(int64(max-min)) 创建一个大的整数，表示生成随机数的范围
	random, err := crand.Int(crand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		rand.Seed(time.Now().UnixNano())
		// 生成一个指定范围内的随机整数
		return rand.Intn(max-min) + min
	}

	return int(random.Int64()) + min
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
