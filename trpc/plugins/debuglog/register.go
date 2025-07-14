package debuglog

import (
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginName = "debuglog"
	pluginType = "tracing"
)

func init() {
	filter.Register("debuglog", ServerFilter(), ClientFilter())
	filter.Register(
		"simpledebuglog", ServerFilter(WithLogFunc(SimpleLogFunc)),
		ClientFilter(WithLogFunc(SimpleLogFunc)),
	)
	filter.Register(
		"pjsondebuglog", ServerFilter(WithLogFunc(PrettyJSONLogFunc)),
		ClientFilter(WithLogFunc(PrettyJSONLogFunc)),
	)
	filter.Register(
		"jsondebuglog", ServerFilter(WithLogFunc(JSONLogFunc)),
		ClientFilter(WithLogFunc(JSONLogFunc)),
	)
}

// Register the plugin on init
func init() {
	plugin.Register(pluginName, &Plugin{})
}
