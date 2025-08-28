package ristretto

import (
	"fmt"
	"testing"

	"github.com/dgraph-io/ristretto/v2"
)

func TestSet(t *testing.T) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M). 跟踪频率的键数量 (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB). 缓存的最大成本 (1GB).
		BufferItems: 64,      // number of keys per Get buffer. 每个Get缓冲的键数量.
	})
	if err != nil {
		panic(err)
	}
	defer cache.Close()

	// set a value with a cost of 1
	cache.Set("key", "value", 1)

	// wait for value to pass through buffers
	cache.Wait()

	// get value from cache
	value, found := cache.Get("key")
	if !found {
		panic("missing value")
	}
	fmt.Println(value)

	// del value from cache
	cache.Del("key")
}
