package masking

import (
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

func init() {
	filter.Register(pluginName, ServerFilter(), nil)
}

func init() {
	plugin.Register(pluginName, &MaskingPlugin{})
}
