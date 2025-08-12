package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
	"trpc.group/trpc-go/tnet/log"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginType = "circuitbreaker"
	pluginName = "hystrix"
)

var cfg map[string]hystrix.CommandConfig

type hystrixPlugin struct{}

// Type ...
func (p *hystrixPlugin) Type() string {
	return pluginType
}

// Setup ...
func (p *hystrixPlugin) Setup(name string, decoder plugin.Decoder) error {
	// 读取插件配置
	cfg = make(map[string]hystrix.CommandConfig)
	err := decoder.Decode(cfg)
	if err != nil {
		log.Errorf("decoder.Decode(%T) err(%s)", cfg, err.Error())
		return err
	}

	// 初始化 hystrix
	hystrix.Configure(cfg)

	// 注册过滤器
	filter.Register(filterName, ServerFilter(), ClientFilter())

	return nil
}
