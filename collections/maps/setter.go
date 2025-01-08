package maps

import (
	"fmt"
	"strings"
)

// SetItems 对嵌套 Map 进行赋值
// paths 参数支持 []string 类型，如 []string{"metadata", "namespace"}
// 或 string 类型（以 '.' 为分隔符），如 "spec.template.spec.containers"
func SetItems(obj map[string]interface{}, paths interface{}, val interface{}) error {
	// 检查 paths 类型
	switch p := paths.(type) {
	case string:
		if err := setItems(obj, strings.Split(p, "."), val); err != nil {
			return err
		}
	case []string:
		if err := setItems(obj, p, val); err != nil {
			return err
		}
	default:
		return ErrInvalidPathType
	}
	return nil
}

func setItems(obj map[string]interface{}, paths []string, val interface{}) error {
	if len(paths) == 0 {
		return fmt.Errorf("paths is empty list")
	}
	if len(paths) == 1 {
		obj[paths[0]] = val
	} else if subMap, ok := obj[paths[0]].(map[string]interface{}); ok {
		return setItems(subMap, paths[1:], val)
	} else {
		return fmt.Errorf("key %s not exists or obj[key] not map[string]interface{} type", paths[0])
	}
	return nil
}

// SetDeepMapValue 根据key设置DeepMap的值.
func SetDeepMapValue(m map[string]interface{}, key string, value interface{}, keyDelim string) {
	key = strings.ToLower(key)
	value = ToCaseInsensitiveValue(value)

	path := strings.Split(key, keyDelim)
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := DeepSearch(m, path[0:len(path)-1])

	// set innermost value
	deepestMap[lastKey] = value
}
