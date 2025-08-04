package fastrand

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

// BenchSink prevents the compiler from optimizing away benchmark loops.
var (
	BenchSink      uint32
	BenchSinkInt32 int32
	BenchSinkInt64 int64
)

// 同个随机种子，随机数求和
func BenchmarkUint32n(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		s := uint32(0)
		for pb.Next() {
			s += Uint32n(1e6)
		}
		atomic.AddUint32(&BenchSink, s)
	})
}

// 不同的随机种子，随机数求和
func BenchmarkRNGUint32n(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var r RNG
		s := uint32(0)
		for pb.Next() {
			s += r.Uint32n(1e6)
		}
		atomic.AddUint32(&BenchSink, s)
	})
}

// 使用全局锁的同一个随机种子，随机数求和
func BenchmarkRNGUint32nWithLock(b *testing.B) {
	var r RNG
	var rMu sync.Mutex
	b.RunParallel(func(pb *testing.PB) {
		s := uint32(0)
		for pb.Next() {
			rMu.Lock()
			s += r.Uint32n(1e6)
			rMu.Unlock()
		}
		atomic.AddUint32(&BenchSink, s)
	})
}

func BenchmarkInt(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		s := Int()
		for pb.Next() {
			s += Int()
		}
		atomic.AddInt32(&BenchSinkInt32, int32(s))
	})
}

func BenchmarkIntn(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		s := Intn(1e2)
		for pb.Next() {
			s += Intn(1e2)
		}
		atomic.AddInt32(&BenchSinkInt32, int32(s))
	})
}

func BenchmarkInt31(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		s := Int31()
		for pb.Next() {
			s += Int31()
		}
		atomic.AddInt32(&BenchSinkInt32, s)
	})
}

func BenchmarkInt31n(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		s := Int31n(1e2)
		for pb.Next() {
			s += Int31n(1e2)
		}
		atomic.AddInt32(&BenchSinkInt32, s)
	})
}

func BenchmarkInt63(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		s := Int63()
		for pb.Next() {
			s += Int63()
		}
		atomic.AddInt64(&BenchSinkInt64, s)
	})
}

func BenchmarkInt63n(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		s := Int63n(1e2)
		for pb.Next() {
			s += Int63n(1e2)
		}
		atomic.AddInt64(&BenchSinkInt64, s)
	})
}

func BenchmarkRNGUint32nArray(b *testing.B) {
	var rr [64]struct {
		r  RNG
		mu sync.Mutex

		// pad prevents from false sharing
		pad [64 - (unsafe.Sizeof(RNG{})+unsafe.Sizeof(sync.Mutex{}))%64]byte
	}
	var n uint32
	b.RunParallel(func(pb *testing.PB) {
		s := uint32(0)
		for pb.Next() {
			idx := atomic.AddUint32(&n, 1)
			r := &rr[idx%uint32(len(rr))]
			r.mu.Lock()
			s += r.r.Uint32n(1e6)
			r.mu.Unlock()
		}
		atomic.AddUint32(&BenchSink, s)
	})
}

func BenchmarkMathRandInt31n(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		s := uint32(0)
		for pb.Next() {
			s += uint32(rand.Int31n(1e6))
		}
		atomic.AddUint32(&BenchSink, s)
	})
}

func BenchmarkMathRandRNGInt31n(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(42))
		s := uint32(0)
		for pb.Next() {
			s += uint32(r.Int31n(1e6))
		}
		atomic.AddUint32(&BenchSink, s)
	})
}

func BenchmarkMathRandRNGInt31nWithLock(b *testing.B) {
	r := rand.New(rand.NewSource(42))
	var rMu sync.Mutex
	b.RunParallel(func(pb *testing.PB) {
		s := uint32(0)
		for pb.Next() {
			rMu.Lock()
			s += uint32(r.Int31n(1e6))
			rMu.Unlock()
		}
		atomic.AddUint32(&BenchSink, s)
	})
}

func BenchmarkMathRandRNGInt31nArray(b *testing.B) {
	var rr [64]struct {
		r  *rand.Rand
		mu sync.Mutex

		// pad prevents from false sharing
		pad [64 - (unsafe.Sizeof(RNG{})+unsafe.Sizeof(sync.Mutex{}))%64]byte
	}
	for i := range rr {
		rr[i].r = rand.New(rand.NewSource(int64(i)))
	}

	var n uint32
	b.RunParallel(func(pb *testing.PB) {
		s := uint32(0)
		for pb.Next() {
			idx := atomic.AddUint32(&n, 1)
			r := &rr[idx%uint32(len(rr))]
			r.mu.Lock()
			s += uint32(r.r.Int31n(1e6))
			r.mu.Unlock()
		}
		atomic.AddUint32(&BenchSink, s)
	})
}
