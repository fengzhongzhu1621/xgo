package run

import (
	"github.com/fengzhongzhu1621/xgo/collections/set"
	"github.com/fengzhongzhu1621/xgo/config"
)

// InitSuperAppCode 初始化超级 app
func InitSuperAppCode() {
	cfg := config.GetGlobalConfig()
	SuperAppCodeSet := set.NewStringSet()
	for _, app_code := range cfg.SuperAppCode {
		SuperAppCodeSet.Add(app_code)
	}
}
