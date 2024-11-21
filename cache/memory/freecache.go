package memory

import (
	"sync"

	"github.com/coocood/freecache"

	log "github.com/sirupsen/logrus"
)

var (
	freeCacheInstance *freecache.Cache
	initOnce          sync.Once
)

func InitCache(capacity int) {
	initOnce.Do(func() {
		freeCacheInstance = freecache.NewCache(capacity * 1024 * 1024)
	})
}

func NewFreeCache() *freecache.Cache {
	if freeCacheInstance == nil {
		log.Fatal("cache not init")
	}
	return freeCacheInstance
}
