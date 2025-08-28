package ccache

import (
	"context"
	"time"

	"github.com/fengzhongzhu1621/xgo/config"
)

func init() {
	globalConfig := config.GetGlobalConfig()

	// 初始化 LRU 缓存
	InitLocalCache(globalConfig)

	// 启动 LUR 缓存异步入库任务
	timeout := 24 * time.Hour
	go AyncRefreshDefaultCache(context.Background(), timeout)
}
