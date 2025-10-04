package config

import (
	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	// ErrConfigNotExist is config not exist error
	ErrConfigNotExist = errors.New("xgo/config: config not exist")

	// ErrProviderNotExist is provider not exist error
	ErrProviderNotExist = errors.New("xgo/config: provider not exist")

	// ErrCodecNotExist is codec not exist error
	ErrCodecNotExist = errors.New("xgo/config: codec not exist")
)

// GetViperString 如果配置中指定key的值为空，则返回默认值.
func GetViperString(key, defaultValue string) string {
	val := viper.GetString(key)
	if val != "" {
		return val
	}
	return defaultValue
}

// UnmarshalKey 将配置解析到结构体.
func UnmarshalKey(key string, obj interface{}) {
	err := viper.UnmarshalKey(key, &obj)
	if err != nil {
		panic(errors.Wrapf(err, "unable to decode '%s' into struct", key))
	}
}

// search 在嵌套的配置数据中递归搜索指定键路径的值
// 参数:
//   unmarshalledData: 已解析的配置数据映射
//   keys: 键路径数组，用于在嵌套结构中定位目标值
// 返回:
//   interface{}: 找到的配置值
//   error: 如果键不存在则返回 ErrConfigNotExist 错误
func search(unmarshalledData map[string]interface{}, keys []string) (interface{}, error) {
	// 如果键路径为空，返回配置不存在错误
	if len(keys) == 0 {
		return nil, ErrConfigNotExist
	}

	// 在当前层级查找第一个键
	key, ok := unmarshalledData[keys[0]]
	if !ok {
		return nil, ErrConfigNotExist
	}

	// 如果只有一个键，直接返回找到的值
	if len(keys) == 1 {
		return key, nil
	}
	
	// 根据值的类型进行递归搜索
	switch key := key.(type) {
	case map[interface{}]interface{}:
		// 处理 interface{} 键的映射，转换为 string 键的映射后递归搜索
		return search(cast.ToStringMap(key), keys[1:])
	case map[string]interface{}:
		// 处理 string 键的映射，直接递归搜索剩余的键路径
		return search(key, keys[1:])
	default:
		// 如果当前值不是映射类型，但还有剩余键路径，说明配置结构不匹配
		return nil, ErrConfigNotExist
	}
}
