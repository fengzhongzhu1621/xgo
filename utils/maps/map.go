package maps

import (
	"strings"
)

// DeepSearchAndCreateMap 将摊平的字典转换成收缩的字典.
func CreateDeepMap(m map[string]interface{}, keyDelim string) map[string]interface{} {
	m2 := map[string]interface{}{}
	// start from the list of keys, and construct the map one value at a time
	for key, value := range m {
		if value == nil {
			// should not happen, since AllKeys() returns only keys holding a value,
			// check just in case anything changes
			continue
		}
		path := strings.Split(key, keyDelim)
		lastKey := strings.ToLower(path[len(path)-1])
		deepestMap := DeepSearch(m2, path[0:len(path)-1])
		// set innermost value
		deepestMap[lastKey] = value
	}
	return m2
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
