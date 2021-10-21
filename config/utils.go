package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// 如果配置中指定key的值为空，则返回默认值
func GetString(key, defaultValue string) string {
	val := viper.GetString(key)
	if val != "" {
		return val
	}
	return defaultValue
}

// 将配置解析到结构体
func UnmarshalKey(key string, obj interface{}) {
	err := viper.UnmarshalKey(key, &obj)
	if err != nil {
		panic(errors.Wrapf(err, "unable to decode '%s' into struct", key))
	}
}
