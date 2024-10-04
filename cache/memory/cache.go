package memory

import (
	"time"

	cache "github.com/fengzhongzhu1621/xgo/cache/common"
	"github.com/fengzhongzhu1621/xgo/cache/memory/backend"
)

func NewCache(
	name string,
	retrieveFunc cache.RetrieveFunc,
	expiration time.Duration,
	randomExtraExpirationFunc cache.RandomExtraExpirationDurationFunc,
	options ...Option,
) Cache {
	cacheBackend := backend.NewMemoryBackend(name, expiration, randomExtraExpirationFunc)
	return NewBaseCache(retrieveFunc, cacheBackend, options...)
}
