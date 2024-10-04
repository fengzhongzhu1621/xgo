package cleaner

import (
	"context"

	cache "github.com/fengzhongzhu1621/xgo/cache/common"
	"github.com/fengzhongzhu1621/xgo/monitor/sentry"
	log "github.com/sirupsen/logrus"
)

// it's a goroutine

// go CacheCleaner.Delete(a)
// go CacheCleaner.Delete([a1, a2, a3])

// then => put into channel

// the consumer:
// a => type => will case some other cache delete
// 例如 delete subject => will delete subject-group / subject-department / subjectpk ....

const defaultCacheCleanerBufferSize = 2000

// CacheDeleter ...
type CacheDeleter interface {
	// 删除缓存中的特定键
	Execute(key cache.Key) error
}

// CacheCleaner 定义了一个缓存清理器，它负责从缓存中删除键值对
type CacheCleaner struct {
	// 缓存清理器的名称
	name string
	// 上下文，用于控制清理器的生命周期
	ctx context.Context
	// 一个带缓冲的通道，用于存储待删除的缓存键
	buffer chan cache.Key
	// 实现了 CacheDeleter 接口的对象，负责实际的删除操作
	deleter CacheDeleter
}

// NewCacheCleaner 创建并初始化一个新的 CacheCleaner 实例
func NewCacheCleaner(name string, deleter CacheDeleter) *CacheCleaner {
	ctx := context.Background()
	return &CacheCleaner{
		name:    name,
		ctx:     ctx,
		buffer:  make(chan cache.Key, defaultCacheCleanerBufferSize),
		deleter: deleter,
	}
}

// Run 启动缓存清理器的运行循环。它会不断从 buffer 中读取键并调用 deleter.Execute 方法进行删除。如果删除失败，会记录错误并通过 Sentry 报告。
func (r *CacheCleaner) Run() {
	log.Infof("running a cache cleaner: %s", r.name)
	var err error
	for {
		select {
		case <-r.ctx.Done():
			// 生命周期结束停止清理
			return
		case d := <-r.buffer:
			// 删除缓存中的特定键
			err = r.deleter.Execute(d)
			if err != nil {
				log.Errorf("delete cache key=%s fail: %s", d.Key(), err)

				// report to sentry
				sentry.ReportToSentry(
					"cache error: delete key fail",
					map[string]interface{}{
						"key":   d.Key(),
						"error": err.Error(),
					},
				)
			}
		}
	}
}

// Delete 将单个键发送到 buffer 中，以便稍后由 Run 方法处理
func (r *CacheCleaner) Delete(key cache.Key) {
	r.buffer <- key
}

// BatchDelete 将多个键批量发送到 buffer 中。目前实现是逐个发送，未来可以考虑支持批量删除的优化。
func (r *CacheCleaner) BatchDelete(keys []cache.Key) {
	// TODO: support batch delete in pipeline or tx?
	for _, key := range keys {
		r.buffer <- key
	}
}
