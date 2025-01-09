package maps

import (
	"fmt"
	"sync"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
)

// ConcurrentMap is like map, but is safe for concurrent use by multiple goroutines.
// NewConcurrentMap create a ConcurrentMap with specific shard count.
// func NewConcurrentMap[K comparable, V any](shardCount int) *ConcurrentMap[K, V]
// func (cm *ConcurrentMap[K, V]) Set(key K, value V)
// func (cm *ConcurrentMap[K, V]) Get(key K) (V, bool)
func TestNewConcurrentMap(t *testing.T) {
	// create a ConcurrentMap whose key type is string, value type is int
	cm := maputil.NewConcurrentMap[string, int](100)

	var wg1 sync.WaitGroup
	wg1.Add(5)

	for i := 0; i < 5; i++ {
		go func(n int) {
			cm.Set(fmt.Sprintf("%d", n), n)
			wg1.Done()
		}(i)
	}
	wg1.Wait()

	var wg2 sync.WaitGroup
	wg2.Add(5)
	for j := 0; j < 5; j++ {
		go func(n int) {
			val, ok := cm.Get(fmt.Sprintf("%d", n))
			fmt.Println(val, ok)
			wg2.Done()
		}(j)
	}
	wg2.Wait()

	// output: (order may change)
	// 1 true
	// 3 true
	// 2 true
	// 0 true
	// 4 true
}

// TestConcurrentMapGetOrSet Returns the existing value for the key if present. Otherwise, it sets and returns the given value.
// func (cm *ConcurrentMap[K, V]) GetOrSet(key K, value V) (actual V, ok bool)
func TestConcurrentMapGetOrSet(t *testing.T) {
	cm := maputil.NewConcurrentMap[string, int](100)

	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(n int) {
			val, ok := cm.GetOrSet(fmt.Sprintf("%d", n), n)
			fmt.Println(val, ok)
			wg.Done()
		}(i)
	}
	wg.Wait()

	// output: (order may change)
	// 1 false
	// 3 false
	// 2 false
	// 0 false
	// 4 false
}

// TestConcurrentMapDelete Delete the value for a key.
// func (cm *ConcurrentMap[K, V]) Delete(key K)
func TestConcurrentMapDelete(t *testing.T) {
	cm := maputil.NewConcurrentMap[string, int](100)

	var wg1 sync.WaitGroup
	wg1.Add(5)

	for i := 0; i < 5; i++ {
		go func(n int) {
			cm.Set(fmt.Sprintf("%d", n), n)
			wg1.Done()
		}(i)
	}
	wg1.Wait()

	var wg2 sync.WaitGroup
	wg2.Add(5)
	for j := 0; j < 5; j++ {
		go func(n int) {
			cm.Delete(fmt.Sprintf("%d", n))
			wg2.Done()
		}(j)
	}

	cm.Delete("unknown")

	wg2.Wait()
}

// TestConcurrentMapGetAndDelete Returns the existing value for the key if present and then delete the value for the key. Otherwise, do nothing, just return false.
// func (cm *ConcurrentMap[K, V]) GetAndDelete(key K) (actual V, ok bool)
func TestConcurrentMapGetAndDelete(t *testing.T) {
	cm := maputil.NewConcurrentMap[string, int](100)

	var wg1 sync.WaitGroup
	wg1.Add(5)

	for i := 0; i < 5; i++ {
		go func(n int) {
			cm.Set(fmt.Sprintf("%d", n), n)
			wg1.Done()
		}(i)
	}
	wg1.Wait()

	var wg2 sync.WaitGroup
	wg2.Add(5)
	for j := 0; j < 5; j++ {
		go func(n int) {
			val, ok := cm.GetAndDelete(fmt.Sprintf("%d", n))
			fmt.Println(val, ok) //n, true

			_, ok = cm.Get(fmt.Sprintf("%d", n))
			fmt.Println(val, ok) //false

			wg2.Done()
		}(j)
	}

	wg2.Wait()

	// 4 true
	// 4 false
	// 3 true
	// 3 false
	// 1 true
	// 1 false
	// 2 true
	// 2 false
	// 0 true
	// 0 false
}

// TestConcurrentMapHas Checks if map has the value for a key.
// func (cm *ConcurrentMap[K, V]) Has(key K) bool
func TestConcurrentMapHas(t *testing.T) {
	cm := maputil.NewConcurrentMap[string, int](100)

	var wg1 sync.WaitGroup
	wg1.Add(5)

	for i := 0; i < 5; i++ {
		go func(n int) {
			cm.Set(fmt.Sprintf("%d", n), n)
			wg1.Done()
		}(i)
	}
	wg1.Wait()

	var wg2 sync.WaitGroup
	wg2.Add(5)
	for j := 0; j < 5; j++ {
		go func(n int) {
			ok := cm.Has(fmt.Sprintf("%d", n))
			fmt.Println(ok) // true
			wg2.Done()
		}(j)
	}
	wg2.Wait()

	// true
	// true
	// true
	// true
	// true

}

// TestConcurrentMapIter Calls iterator sequentially for each key and value present in each of the shards in the map. If iterator returns false, range stops the iteration.
// func (cm *ConcurrentMap[K, V]) Range(iterator func(key K, value V) bool)
func TestConcurrentMapIter(t *testing.T) {
	cm := maputil.NewConcurrentMap[string, int](100)

	var wg1 sync.WaitGroup
	wg1.Add(5)

	for i := 0; i < 5; i++ {
		go func(n int) {
			cm.Set(fmt.Sprintf("%d", n), n)
			wg1.Done()
		}(i)
	}
	wg1.Wait()

	cm.Range(func(key string, value int) bool {
		fmt.Println(value)
		return true
	})

	// 4
	// 3
	// 2
	// 1
	// 0
}
