package maps

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidPathType = errors.New("paths's type must one of (string, []string)")

// GetItems 根据路径从字典中取值
func GetItems(obj map[string]interface{}, paths interface{}) (interface{}, error) {
	switch p := paths.(type) {
	case string:
		return getItems(obj, strings.Split(p, "."))
	case []string:
		return getItems(obj, p)
	default:
		return nil, ErrInvalidPathType
	}
}

func getItems(obj map[string]interface{}, paths []string) (interface{}, error) {
	if len(paths) == 0 {
		return nil, errors.New("paths is empty list")
	}
	ret, exists := obj[paths[0]]
	if !exists {
		return nil, fmt.Errorf("key %s not exist", paths[0])
	}
	if len(paths) == 1 {
		return ret, nil
	} else if subMap, ok := obj[paths[0]].(map[string]interface{}); ok {
		// 递归获取子字典的值
		return getItems(subMap, paths[1:])
	}
	return nil, fmt.Errorf("key %s, val not map[string]interface{} type", paths[0])
}

// Get 若指定值不存在，则返回默认值
func Get(obj map[string]interface{}, paths interface{}, defVal interface{}) interface{} {
	ret, err := GetItems(obj, paths)
	if err != nil {
		return defVal
	}
	return ret
}

// GetBool 获取 Bool 类型快捷方法，默认值为 false
func GetBool(obj map[string]interface{}, paths interface{}) bool {
	return Get(obj, paths, false).(bool)
}

// GetInt64 获取 int64 类型快捷方法，默认值为 int64(0)
func GetInt64(obj map[string]interface{}, paths interface{}) int64 {
	return Get(obj, paths, int64(0)).(int64)
}

// GetStr 获取 string 类型快捷方法，默认值为 ""
func GetStr(obj map[string]interface{}, paths interface{}) string {
	return Get(obj, paths, "").(string)
}

// GetList 获取 []interface{} 类型快捷方法，默认值为 []interface{}{}
func GetList(obj map[string]interface{}, paths interface{}) []interface{} {
	return Get(obj, paths, []interface{}{}).([]interface{})
}

// GetMap 获取 map[string]interface{} 类型快捷方法，默认值为 map[string]interface{}
func GetMap(obj map[string]interface{}, paths interface{}) map[string]interface{} {
	return Get(obj, paths, map[string]interface{}{}).(map[string]interface{})
}
