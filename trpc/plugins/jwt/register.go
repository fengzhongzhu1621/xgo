package jwt

import "trpc.group/trpc-go/trpc-go/plugin"

const (
	pluginName = "jwt"
	pluginType = "auth"
)

func init() {
	plugin.Register(pluginName, &pluginImp{})
}
