package mapstructure

// 将 Go 结构体反向解码为 map[string]interface{}

import "github.com/mitchellh/mapstructure"

// A DecoderConfigOption can be passed to viper.Unmarshal to configure
// mapstructure.DecoderConfig options .
type DecoderConfigOption func(*mapstructure.DecoderConfig)

// DecodeHook returns a DecoderConfigOption which overrides the default
// DecoderConfig.DecodeHook value, the default is:
//
//  mapstructure.ComposeDecodeHookFunc(
//		mapstructure.StringToTimeDurationHookFunc(),
//		mapstructure.StringToSliceHookFunc(","),
//	)
// 给 mapstructure.DecoderConfig 添加 hook函数
func DecodeHook(hook mapstructure.DecodeHookFunc) DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.DecodeHook = hook
	}
}
