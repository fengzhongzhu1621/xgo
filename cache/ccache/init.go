package ccache

import (
	"github.com/fengzhongzhu1621/xgo/config"
)

func init() {
	globalConfig := config.GetGlobalConfig()
	InitLocalCache(globalConfig)
}
