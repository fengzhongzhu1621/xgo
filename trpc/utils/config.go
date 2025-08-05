package utils

import (
	"gopkg.in/yaml.v3"
	"trpc.group/trpc-go/trpc-go"
)

// ParsePluginConf 解析指定插件的配置
func ParsePluginConf(conf string, pluginType, pluginName string) *yaml.Node {
	// 解析插件的yaml配置
	cfg := trpc.Config{}
	if err := yaml.Unmarshal([]byte(conf), &cfg); err != nil {
		return nil
	}
	// 返回 yaml.Node
	validCfg := cfg.Plugins[pluginType][pluginName]

	return &validCfg
}
