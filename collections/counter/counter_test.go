package counter

// go test -bench=. -benchmem
//
// BenchmarkMutexCounter-10     	 3022845	       342.6 ns/op	      24 B/op	       1 allocs/op
// BenchmarkAtomicCounter-10    	 3925063	       299.8 ns/op	      24 B/op	       1 allocs/op

import (
	"sync"
	"testing"
)

func BenchmarkMutexCounter(b *testing.B) {
	var wg sync.WaitGroup
	counter := MutexCounter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			counter.Increment()
			wg.Done()
		}()
	}
	wg.Wait()
	_ = counter.Value() // 防止编译器优化掉
}

func BenchmarkAtomicCounter(b *testing.B) {
	var wg sync.WaitGroup
	counter := AtomicCounter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			counter.Increment()
			wg.Done()
		}()
	}
	wg.Wait()
	_ = counter.Value() // 防止编译器优化掉
}
