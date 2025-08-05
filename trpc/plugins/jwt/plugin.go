package jwt

import (
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/network/nethttp/auth/jwt"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

type pluginImp struct{}

func (p *pluginImp) Type() string {
	return pluginType
}

// Setup 插件实例初始化
func (p *pluginImp) Setup(name string, configDec plugin.Decoder) error {
	// 配置解析
	cfg := Config{}
	if err := configDec.Decode(&cfg); err != nil {
		return err
	}
	if cfg.Secret == "" {
		return fmt.Errorf("JWT secret not be empty")
	}

	// 未设置-默认为1小时过期
	expired := time.Hour
	if cfg.Expired > 0 {
		expired = time.Duration(cfg.Expired) * time.Second
	}

	// slice 转换为 Set
	excludePathSet := make(map[string]bool)
	for _, s := range cfg.ExcludePaths {
		excludePathSet[s] = true
	}

	// 设置默认 signer
	SetDefaultSigner(jwt.NewJwtSign([]byte(cfg.Secret), expired, cfg.Issuer))

	// 注册过滤器
	filter.Register(name, ServerFilter(WithExcludePathSet(excludePathSet)), nil)

	return nil
}
