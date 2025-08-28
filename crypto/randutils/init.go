package randutils

import (
	"math/rand"
	"time"
)

var globalRNG *rand.Rand

func init() {
	rand.Seed(time.Now().UnixNano())
}

func init() {
	// 建议尽量避免使用全局生成器，除非确实有必要。
	// rand.NewSource 和 rand.New 的开销很小，通常不需要担心性能问题。
	// 但在极高性能要求的场景下，可以进行基准测试（benchmarking）以确保没有性能瓶颈。
	seed := rand.NewSource(time.Now().UnixNano())
	globalRNG = rand.New(seed)
}
