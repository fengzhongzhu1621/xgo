package memory

import (
	"time"

	"github.com/fengzhongzhu1621/xgo/cache"
	"github.com/fengzhongzhu1621/xgo/cache/memory/backend"
)

func NewCache(
	name string,
	retrieveFunc RetrieveFunc,
	expiration time.Duration,
	randomExtraExpirationFunc cache.RandomExtraExpirationDurationFunc,
	options ...Option,
) Cache {
	be := backend.NewMemoryBackend(name, expiration, randomExtraExpirationFunc)
	return NewBaseCache(retrieveFunc, be, options...)
}

// NewMockCache create a memory cache for mock
func NewMockCache(retrieveFunc RetrieveFunc) Cache {
	be := backend.NewMemoryBackend("mockCache", 5*time.Minute, nil)

	return NewBaseCache(retrieveFunc, be)
}
