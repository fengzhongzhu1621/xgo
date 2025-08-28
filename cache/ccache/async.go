package ccache

import (
	"context"
	"time"

	redisClient "github.com/fengzhongzhu1621/xgo/db/redis/client"
	"github.com/go-redis/redis/v8"
)

var (
	MAX_PIPE_NUM                            = 500
	defaultCacheChan chan map[string]string = make(chan map[string]string, 4096)
)

// SetToDefaultCache 异步缓存字典
func SetToDefaultCache(kvmap map[string]string) {
	select {
	// 将一个键值对映射 kvmap 发送到一个全局的缓存通道 defaultCacheChan 中
	case defaultCacheChan <- kvmap:
	// 如果在100毫秒内无法发送（即通道已满），则放弃发送操作
	case <-time.After(time.Duration(100) * time.Millisecond):
	}
}

// GetFromDefaultCache 从缓存中获取数据
func GetFromDefaultCache(kvmap map[string]string) map[string]string {
	var (
		cli         = redisClient.GetDefaultRedisClient()
		pipeline    = cli.Pipeline()
		missing     = false
		pipe_result = make(map[string]*redis.StringCmd)
	)

	// 输入参数为空
	if kvmap == nil {
		return nil
	}

	// 获得 LRU 缓存
	cache := GetDefaultCache()

	for k := range kvmap {
		// 从 LRU 中获取缓存的结果
		item := cache.Get(k)
		if item == nil {
			// 如果没有命中，则尝试从 redis 中获取结果
			missing = true
			pipe_result[k] = pipeline.Get(context.Background(), k)
		} else {
			if item.Expired() {
				// 本地缓存存在但是已经过期，则主动删除
				kvmap[k] = item.Value().(string)
				cache.Delete(k)
			} else {
				// 本地缓存存在未过期
				kvmap[k] = item.Value().(string)
			}
		}
	}

	// 本地缓存没有命中，从 redis 获取业务名
	if missing {
		_, _ = pipeline.Exec(context.Background())
		for key, r := range pipe_result {
			val, err := r.Result()
			if err == nil {
				kvmap[key] = val
			}
		}

	}

	return kvmap
}

// AyncRefreshDefaultCache 启动 LUR 缓存异步入库任务
func AyncRefreshDefaultCache(ctx context.Context, timeout time.Duration) {
	cli := redisClient.GetDefaultRedisClient()

	for {
		select {
		case <-ctx.Done():
			return
		case dataMap := <-defaultCacheChan:
			// 从队列获取需要缓存的数据
			if dataMap == nil || len(dataMap) == 0 {
			} else {
				index := 0
				pipeline := cli.Pipeline()
				// 获得 LRU 缓存
				cache := GetDefaultCache()
				for k, v := range dataMap {
					// 缓存数据
					cache.Set(k, v, timeout)
					pipeline.Set(ctx, k, v, timeout)
					index++
					// 如果缓存的数据量比较多，则分批次缓存
					if index%MAX_PIPE_NUM == 0 {
						pipeline.Exec(ctx)
					}
				}

				if index%MAX_PIPE_NUM != 0 {
					pipeline.Exec(ctx)
				}
			}
		}
	}
}
