package transinfoblocker

import (
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

func init() {
	plugin.Register(pluginName, &TransinfoBlocker{})
}

func init() {
	filter.Register("transinfo-blocker", filter.NoopServerFilter, ClientFilter)
}
