package validation

import (
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginName = "validation"
	pluginType = "auth"
)

func init() {
	// 注册拦截器
	filter.Register(
		pluginName,
		ServerFilterWithOptions(defaultOptions),
		ClientFilterWithOptions(defaultOptions),
	)
}

func init() {
	// 注册插件
	plugin.Register(pluginName, &ValidationPlugin{})
}
