package hystrix

import "trpc.group/trpc-go/trpc-go/plugin"

func init() {
	plugin.Register(pluginName, &hystrixPlugin{})
}
