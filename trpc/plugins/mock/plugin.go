package mock

import (
	"encoding/base64"
	"time"

	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const (
	pluginName = "mock"
	pluginType = "tracing"
)

// Plugin mock trpc plugin implementation.
type Plugin struct{}

// Type mock trpc plugin type.
func (p *Plugin) Type() string {
	return pluginType
}

// Setup mock instance initialization.
func (p *Plugin) Setup(name string, configDec plugin.Decoder) error {
	// 读取插件配置
	conf := Config{}
	if err := configDec.Decode(&conf); err != nil {
		return err
	}

	var opt []Option
	for _, mock := range conf {
		mock.delay = time.Millisecond * time.Duration(mock.Delay)
		// 格式化响应body
		if mock.Serialization != codec.SerializationTypeJSON {
			// When the serialization method is not json, use base64 to decode.
			decoded, err := base64.StdEncoding.DecodeString(mock.Body)
			if err != nil {
				return err
			}
			mock.data = decoded
		} else {
			mock.data = []byte(mock.Body)
		}
		opt = append(opt, WithMock(mock))
	}

	filter.Register(pluginName, nil, ClientFilter(opt...))
	return nil
}
