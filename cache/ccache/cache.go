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

		// MaxSize(409600) 设置了缓存的最大条目数。在这个例子中，缓存最多可以存储409600个条目。
		// 当缓存中的条目数超过这个限制时，将会根据一定的策略（如最近最少使用LRU）移除一些条目以腾出空间
		// 设置一个很大的MaxSize值会消耗更多的内存。务必根据你的应用实际需求和可用资源来合理设置这个参数
		//
		// ItemsToPrune(4096): 指定了在进行缓存清理时要移除的最少条目数。
		// 在这个例子中，当需要清理缓存时，至少会移除4096个条目。
		// 这有助于确保即使在缓存接近其最大容量时，也能保持一定的性能和可用性。
		// 虽然ItemsToPrune设置了最小移除数量，但具体的清理时机和策略还是由ccache库内部决定的。
		localCertCache = ccache.New(ccache.Configure().MaxSize(409600).ItemsToPrune(4096))
	})
}

func GetLocalCertCache() *ccache.Cache {
	return localCertCache
}
