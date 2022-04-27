package utils

import (
	"strings"

	"xgo/cast"
)

// ToCaseInsensitiveValue 将字典的key转换为小写，返回新的字典
// toCaseInsensitiveValue checks if the value is a  map;
// if so, create a copy and lower-case the keys recursively.
func ToCaseInsensitiveValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		value = copyAndInsensitiviseMap(cast.ToStringMap(v))
	case map[string]interface{}:
		value = copyAndInsensitiviseMap(v)
	}

	return value
}

// copyAndInsensitiviseMap behaves like insensitiviseMap, but creates a copy of
// any map it makes case insensitive.
func copyAndInsensitiviseMap(m map[string]interface{}) map[string]interface{} {
	nm := make(map[string]interface{})

	for key, val := range m {
		// 字典的key转换为小写
		lkey := strings.ToLower(key)
		// 判断value是否也是字典
		switch v := val.(type) {
		case map[interface{}]interface{}:
			nm[lkey] = copyAndInsensitiviseMap(cast.ToStringMap(v))
		case map[string]interface{}:
			nm[lkey] = copyAndInsensitiviseMap(v)
		default:
			nm[lkey] = v
		}
	}

	return nm
}

// InsensitiviseMap 将字典的key转换为小写，返回原字典
func InsensitiviseMap(m map[string]interface{}) {
	for key, val := range m {
		switch val.(type) {
		case map[interface{}]interface{}:
			// nested map: cast and recursively insensitivise
			val = cast.ToStringMap(val)
			InsensitiviseMap(val.(map[string]interface{}))
		case map[string]interface{}:
			// nested map: recursively insensitivise
			InsensitiviseMap(val.(map[string]interface{}))
		}

		lower := strings.ToLower(key)
		if key != lower {
			// remove old key (not lower-cased)
			delete(m, key)
		}
		// update map
		m[lower] = val
	}
}
