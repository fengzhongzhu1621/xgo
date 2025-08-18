package transinfoblocker

import (
	"errors"

	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginName = "transinfo-blocker"
	pluginType = "security"
)

// TransinfoBlocker 安全插件
type TransinfoBlocker struct{}

// Type transinfo-blocker 插件类型
func (t *TransinfoBlocker) Type() string {
	return pluginType
}

// Setup transinfo-blocker 实例初始化
func (t *TransinfoBlocker) Setup(name string, configDec plugin.Decoder) error {
	if configDec == nil {
		return errors.New("transinfo-blocker configDec nil")
	}

	// 读取插件配置
	cfg = &Config{}
	if err := configDec.Decode(cfg); err != nil {
		return err
	}

	// 格式化配置
	if err := cfg.Default.parseBlockSet(); err != nil {
		return err
	}

	for _, v := range cfg.RPCNameCfg {
		if err := v.parseBlockSet(); err != nil {
			return err
		}
	}
	return nil
}
