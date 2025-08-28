package referer

import (
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginName = "referer"
	pluginType = "auth"
)

func init() {
	plugin.Register(pluginName, &Plugin{})
}

func init() {
	filter.Register(pluginName, ServerFilter(), nil)
}
