package cacheimpls

import (
	cache "github.com/fengzhongzhu1621/xgo/cache/common"
	"go.uber.org/multierr"
)

type testCacheDeleter struct{}

// Execute 模拟删除缓存操作
func (d testCacheDeleter) Execute(key cache.Key) (err error) {
	err = multierr.Combine(
		TestRedisCache.Delete(key),
	)
	return
}
