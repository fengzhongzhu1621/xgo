package ccache

import (
	"errors"
	"sync"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/karlseguin/ccache/v2"
)

var (
	once           sync.Once
	localCertCache *ccache.Cache
)

func InitLocalCache(conf *config.Config) {
	once.Do(func() {
		if conf == nil {
			panic(errors.New("cache config error"))
		}

		localCertCache = ccache.New(ccache.Configure().MaxSize(409600).ItemsToPrune(4096))
	})
}

func GetLocalCertCache() *ccache.Cache {
	return localCertCache
}
