package validation

import (
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

// ValidationPlugin implements the trpc validation plugin.
type ValidationPlugin struct{} //nolint:revive

// Type validation trpc plugin type.
func (p *ValidationPlugin) Type() string {
	return pluginType
}

// Setup initializes the validation plugin instance.
func (p *ValidationPlugin) Setup(name string, configDec plugin.Decoder) error {
	o := defaultOptions

	// 解析配置并覆盖默认配置
	// When the configuration is not empty,
	// the default values are overridden during configuration parsing.
	if configDec != nil {
		if err := configDec.Decode(&o); err != nil {
			// 插件配置解析错误
			return err
		}
	}

	// 自动注册拦截器
	filter.Register(pluginName, ServerFilterWithOptions(o), ClientFilterWithOptions(o))
	return nil
}
