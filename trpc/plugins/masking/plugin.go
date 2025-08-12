package masking

import (
	"errors"

	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginName = "masking"
	pluginType = "auth"
)

// MaskingPlugin masking trpc 插件实现
type MaskingPlugin struct{}

// Type masking trpc插件类型
func (p *MaskingPlugin) Type() string {
	return pluginType
}

// Setup masking实例初始化
func (p *MaskingPlugin) Setup(_ string, configDec plugin.Decoder) error {
	// 配置解析
	if configDec == nil {
		return errors.New("masking decoder empty")
	}
	conf := make(map[string][]string)
	err := configDec.Decode(&conf)
	if err != nil {
		return err
	}

	sf := ServerFilter()

	filter.Register(pluginName, sf, nil)
	return nil
}
