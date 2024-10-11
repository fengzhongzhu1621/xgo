package randutils

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	sources     = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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

// 支持多线程安全.
func (s *source) Int63() int64 {
	s.mu.Lock()
	n := s.src.Int63()
	s.mu.Unlock()
	return n
}

// 支持多线程安全.
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

func RandAuthToken() string {
	buf := make([]byte, 32)
	_, err := crand.Read(buf)
	if err != nil {
		return RandString2(64)
	}

	return fmt.Sprintf("%x", buf)
}
