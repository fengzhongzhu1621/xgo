package randutils

import (
	crand "crypto/rand"
	"math/rand"
	"time"
)

// Random 随机数据
func Random(size int) []byte {
	buf := make([]byte, size)

	if _, err := crand.Read(buf); err != nil {
		rand.Seed(time.Now().UnixNano())
		rand.Read(buf)
	}
	return buf
}
