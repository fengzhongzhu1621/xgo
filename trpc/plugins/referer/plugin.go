package referer

import (
	"errors"

	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

// Plugin  插件实现
type Plugin struct{}

// Type Plugin trpc插件类型
func (p *Plugin) Type() string {
	return pluginType
}

// Setup Referer实例初始化
func (p *Plugin) Setup(name string, configDec plugin.Decoder) error {
	// 配置解析
	if configDec == nil {
		return errors.New("referer writer decoder empty")
	}

	// 获取yaml配置
	conf := make(map[string][]string)
	if err := configDec.Decode(&conf); err != nil {
		return err
	}

	var opt []Option
	for methodName, allowDomain := range conf {
		opt = append(opt, WithRefererDomain(methodName, allowDomain...))
	}

	filter.Register(pluginName, ServerFilter(opt...), nil)
	return nil
}
