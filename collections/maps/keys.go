package maps

import "strings"

// MapKeysValueIsStruct 获得字典的所有key值.
func MapKeysValueIsStruct(m map[string]struct{}) []string {
	s := make([]string, 0, len(m))
	i := 0
	for k := range m {
		s[i] = k
		i++
	}
	return s
}

// MapKeys 获得字典的所有key值.
func MapKeys(m map[string]interface{}) []string {
	s := make([]string, 0, len(m))
	i := 0
	for k := range m {
		s[i] = k
		i++
	}
	return s
}

// ExistsKey 判断 key 是否存在于 map 中
func ExistsKey(obj map[string]interface{}, key string) bool {
	_, ok := obj[key]
	return ok
}

// KeyExists 判断字段中是否存在指定key，key不区分大小写，返回源key.
func KeyExists(k string, m map[string]interface{}) string {
	lk := strings.ToLower(k)
	for mk := range m {
		lmk := strings.ToLower(mk)
		if lmk == lk {
			return mk
		}
	}
	return ""
}
