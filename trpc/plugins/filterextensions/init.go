package filterextensions

import (
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	// PluginType plugin 的类型
	PluginType = "filter_extensions"
	// PluginName plugin 的名字
	PluginName = "method_filters"
)

func init() {
	plugin.Register(PluginName, &serviceMethodFiltersPlugin{})
}
