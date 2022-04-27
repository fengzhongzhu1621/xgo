// Package fastrand implements fast pesudorandom number generator
// that should scale well on multi-CPU systems.
//
// Use crypto/rand instead of this package for generating
// cryptographically secure random numbers.
package fastrand

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	// 初始化时自动提供一个随机种子
	rand.Seed(time.Now().UnixNano())
}

// 存放可被重复使用的值的容器，缓存随机数.
var rngPool sync.Pool

// 根据时间戳获取无符号32位整数.
func getRandomUint32() uint32 {
	// 纳秒时间戳
	x := time.Now().UnixNano()
	return uint32((x >> 32) ^ x)
}

// RNG is a pseudorandom number generator.
//
// It is unsafe to call RNG methods from concurrent goroutines.
type RNG struct {
	x uint32
}

// Uint32 returns pseudorandom uint32.
//
// It is unsafe to call this method from concurrent goroutines.
// 随机范围 [0, MaxInt64].
func (r *RNG) Uint32() uint32 {
	for r.x == 0 {
		// 获取无符号32位整数作为随机种子
		r.x = getRandomUint32()
	}

	// 采用 XorShift 计算无符号32位随机数
	// XorShift随机数生成器，也称为移位寄存器生成器，是George Marsaglia发现的一类伪随机数生成器。
	// 它是线性反馈移位寄存器（LFSR）的子集，它们允许在软件中进行特别有效的实现，而无需使用过于稀疏的多项式。
	// See https://en.wikipedia.org/wiki/Xorshift
	x := r.x
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	r.x = x
	return x
}

// Uint32n returns pseudorandom uint32 in the range [0..maxN).
//
// It is unsafe to call this method from concurrent goroutines.
// 并发调用可能导致随机数重复.
func (r *RNG) Uint32n(maxN uint32) uint32 {
	// 获得随机数，随机范围是 [0, MaxInt64]
	x := r.Uint32()
	// See http://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
	return uint32((uint64(x) * uint64(maxN)) >> 32)
}

// Seed sets the r state to n.
// 重置随机种子，并发安全.
func (r *RNG) Seed(seed uint32) bool {
	oldSeed := atomic.LoadUint32(&r.x)
	if oldSeed == seed {
		return false
	}
	return atomic.CompareAndSwapUint32(&r.x, oldSeed, seed)
}

// Uint32 returns pseudorandom uint32，并发安全。
//
// It is safe calling this function from concurrent goroutines.
func Uint32() uint32 {
	// 从缓存中获取上一个随机数
	// 如果有并发出现，v必然为nil，会生成一个具有新的随机种子的随机数
	v := rngPool.Get()
	if v == nil {
		v = &RNG{}
	}
	r := v.(*RNG)
	// 获得一个 uint32 随机数，可以并发获取，因为获取上一个随机数时是并发安全的
	x := r.Uint32()
	// 缓存最新的一个随机数，解决 math/rand 中的全局锁问题
	rngPool.Put(r)
	return x
}

// Uint32n returns pseudorandom uint32 in the range [0..maxN)，并发安全。
//
// It is safe calling this function from concurrent goroutines.
func Uint32n(maxN uint32) uint32 {
	x := Uint32()
	// See http://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
	return uint32((uint64(x) * uint64(maxN)) >> 32)
}

type FastRNG struct {
	seed    int64
	rngPool *sync.Pool
}

func NewFastRand() *FastRNG {
	scalableRand := &FastRNG{
		rngPool: &sync.Pool{},
	}
	return scalableRand
}

// 并发安全的设置新的随机种子.
func (r *FastRNG) Seed(newSeed int64) bool {
	// 获得旧的随机种子
	oldSeed := atomic.LoadInt64(&r.seed)
	if oldSeed == newSeed {
		return false
	}
	// 先比较变量的值是否等于给定旧值，等于旧值的情况下才赋予新值，最后返回新值是否设置成功
	return atomic.CompareAndSwapInt64(&r.seed, oldSeed, newSeed)
}

// 创建一个新的随机数发生器.
func (r *FastRNG) NewRand() *rand.Rand {
	// 生成新的随机种子，确保并发安全
	newSeed := time.Now().UnixNano()
	r.Seed(newSeed)
	// 不管随机种子有没有更新成功，都会用新的种子生成随机数发生器
	newRand := rand.New(rand.NewSource(newSeed))
	return newRand
}

var globalFastRand = NewFastRand()

func Int() int {
	// 从缓存中获取上一个随机数
	// 如果有并发出现，v必然为nil，会生成一个具有新的随机种子的发生器
	v := globalFastRand.rngPool.Get()
	if v == nil {
		v = globalFastRand.NewRand()
	}
	r := v.(*rand.Rand)
	// 获得一个随机数，可以并发获取，因为获取上一个随机数时是并发安全的
	x := r.Int()
	// 缓存最新的随机数发生器，解决 math/rand 中的全局锁问题
	globalFastRand.rngPool.Put(v)

	return x
}

func Intn(n int) int {
	// 从缓存中获取上一个随机数
	// 如果有并发出现，v必然为nil，会生成一个具有新的随机种子的发生器
	v := globalFastRand.rngPool.Get()
	if v == nil {
		v = globalFastRand.NewRand()
	}
	r := v.(*rand.Rand)
	// 获得一个随机数，可以并发获取，因为获取上一个随机数时是并发安全的
	x := r.Intn(n)
	// 缓存最新的随机数发生器，解决 math/rand 中的全局锁问题
	globalFastRand.rngPool.Put(v)

	return x
}

func Int31() int32 {
	// 从缓存中获取上一个随机数
	// 如果有并发出现，v必然为nil，会生成一个具有新的随机种子的发生器
	v := globalFastRand.rngPool.Get()
	if v == nil {
		v = globalFastRand.NewRand()
	}
	r := v.(*rand.Rand)
	// 获得一个随机数，可以并发获取，因为获取上一个随机数时是并发安全的
	x := r.Int31()
	// 缓存最新的随机数发生器，解决 math/rand 中的全局锁问题
	globalFastRand.rngPool.Put(v)

	return x
}

func Int31n(n int32) int32 {
	// 从缓存中获取上一个随机数
	// 如果有并发出现，v必然为nil，会生成一个具有新的随机种子的发生器
	v := globalFastRand.rngPool.Get()
	if v == nil {
		v = globalFastRand.NewRand()
	}
	r := v.(*rand.Rand)
	// 获得一个随机数，可以并发获取，因为获取上一个随机数时是并发安全的
	x := r.Int31n(n)
	// 缓存最新的随机数发生器，解决 math/rand 中的全局锁问题
	globalFastRand.rngPool.Put(v)

	return x
}

func Int63() int64 {
	// 从缓存中获取上一个随机数
	// 如果有并发出现，v必然为nil，会生成一个具有新的随机种子的发生器
	v := globalFastRand.rngPool.Get()
	if v == nil {
		v = globalFastRand.NewRand()
	}
	r := v.(*rand.Rand)
	// 获得一个随机数，可以并发获取，因为获取上一个随机数时是并发安全的
	x := r.Int63()
	// 缓存最新的随机数发生器，解决 math/rand 中的全局锁问题
	globalFastRand.rngPool.Put(v)

	return x
}

func Int63n(n int64) int64 {
	// 从缓存中获取上一个随机数
	// 如果有并发出现，v必然为nil，会生成一个具有新的随机种子的发生器
	v := globalFastRand.rngPool.Get()
	if v == nil {
		v = globalFastRand.NewRand()
	}
	r := v.(*rand.Rand)
	// 获得一个随机数，可以并发获取，因为获取上一个随机数时是并发安全的
	x := r.Int63n(n)
	// 缓存最新的随机数发生器，解决 math/rand 中的全局锁问题
	globalFastRand.rngPool.Put(v)

	return x
}
