package randutils

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"time"
)

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
