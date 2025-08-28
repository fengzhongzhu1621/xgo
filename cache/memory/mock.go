package memory

import (
	"time"

	cache "github.com/fengzhongzhu1621/xgo/cache/common"
	"github.com/fengzhongzhu1621/xgo/cache/memory/backend"
)

// NewMockCache create a memory cache for mock
func NewMockCache(retrieveFunc cache.RetrieveFunc) Cache {
	cacheBackend := backend.NewMemoryBackend("mockCache", 5*time.Minute, nil)

	return NewBaseCache(false, retrieveFunc, cacheBackend)
}
