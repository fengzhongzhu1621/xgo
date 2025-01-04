package lru

import (
	"testing"

	"github.com/duke-git/lancet/v2/algorithm"
	"github.com/stretchr/testify/assert"
)

// TestLRUCache LRU缓存实现了具有最近最少使用（LRU）策略的内存缓存。
// func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V]
// func (l *LRUCache[K, V]) Get(key K) (V, bool)
// func (l *LRUCache[K, V]) Put(key K, value V)
// func (l *LRUCache[K, V]) Delete(key K) bool
// func (l *LRUCache[K, V]) Len() int
func TestLRUCache(t *testing.T) {
	cache := algorithm.NewLRUCache[string, int](2)

	cache.Put("a", 1)
	cache.Put("b", 2)

	result1, ok1 := cache.Get("a")
	result2, ok2 := cache.Get("b")
	result3, ok3 := cache.Get("c")

	assert.Equal(t, 1, result1)
	assert.Equal(t, true, ok1)

	assert.Equal(t, 2, result2)
	assert.Equal(t, true, ok2)

	assert.Equal(t, 0, result3)
	assert.Equal(t, false, ok3)

	assert.Equal(t, 2, cache.Len())
	ok := cache.Delete("b")
	assert.Equal(t, true, ok)
}
