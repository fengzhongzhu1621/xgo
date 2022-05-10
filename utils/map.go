package utils

import (
	"fmt"
	"strconv"
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

// InsensitiviseMap 将字典的key转换为小写，返回原字典.
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

// FlattenAndMergeMap recursively flattens the given map into a new map
// Code is based on the function with the same name in tha main package.
// shadow: 摊平后的结果
// m：需要摊平的数组.
func FlattenAndMergeMap(shadow map[string]interface{}, m map[string]interface{}, prefix string, delimiter string) map[string]interface{} {
	if shadow != nil && prefix != "" && shadow[prefix] != nil {
		// prefix is shadowed => nothing more to flatten
		return shadow
	}
	if shadow == nil {
		shadow = make(map[string]interface{})
	}

	var m2 map[string]interface{}
	if prefix != "" {
		prefix += delimiter
	}
	for k, val := range m {
		// 获得完整的key
		fullKey := prefix + k
		switch val.(type) {
		case map[string]interface{}:
			m2 = val.(map[string]interface{})
		case map[interface{}]interface{}:
			m2 = cast.ToStringMap(val)
		default:
			// immediate value
			shadow[strings.ToLower(fullKey)] = val

			continue
		}
		// recursively merge to shadow map
		shadow = FlattenAndMergeMap(shadow, m2, fullKey, delimiter)
	}
	return shadow
}

// SearchMap recursively searches for a value for path in source map.
// Returns nil if not found.
// Note: This assumes that the path entries and map keys are lower cased.
// 在m中搜索路径path
func SearchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		// Fast path
		if len(path) == 1 {
			return next
		}

		// Nested case
		switch next.(type) {
		case map[interface{}]interface{}:
			return SearchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return SearchMap(next.(map[string]interface{}), path[1:])
		default:
			// got a value but nested key expected, return "nil" for not found
			return nil
		}
	}
	return nil
}

// DeepSearch THIS CODE IS COPIED HERE: IT SHOULD NOT BE MODIFIED
// AT SOME POINT IT WILL BE MOVED TO A COMMON PLACE
// deepSearch scans deep maps, following the key indexes listed in the
// sequence "path".
// The last value is expected to be another map, and is returned.
//
// In case intermediate keys do not exist, or map to a non-map value,
// a new map is created and inserted, and the search continues from there:
// the initial map "m" may be modified!
// 深度遍历并构造字典，根据path构造字典，并返回指定路径上的最后一个字典，与 python 中的 setdefault 方法类似.
func DeepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	// 遍历 path 数组
	for _, k := range path {
		// 判断是否在字典m中存在
		m2, ok := m[k]
		if !ok {
			// 如果path不在字典中，将其放到字典m中
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]interface{})
			m[k] = m3 // 修改输入参数的值
			m = m3
			continue
		}
		// 如果path在字典中
		m3, ok := m2.(map[string]interface{})
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}

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

// MergeFlatMap merges the given maps, excluding values of the second map
// shadowed by values from the first map.
// 合并字典的key， value不可并；如果key重复，不合并到源字典.
func MergeFlatMap(shadow map[string]bool, src map[string]interface{}, keyDelim string) map[string]bool {
	// scan keys
outer:
	for k := range src {
		path := strings.Split(k, keyDelim)
		// scan intermediate paths
		var parentKey string
		for i := 1; i < len(path); i++ {
			parentKey = strings.Join(path[0:i], keyDelim)
			if shadow[parentKey] {
				// path is shadowed, continue
				// 如果key已经存在，则检测下一个key
				continue outer
			}
		}
		// add key 标记key已存在
		shadow[strings.ToLower(k)] = true
	}
	return shadow
}

// mergeMaps merges two maps. The `itgt` parameter is for handling go-yaml's
// insistence on parsing nested structures as `map[interface{}]interface{}`
// instead of using a `string` as the key for nest structures beyond one level
// deep. Both map types are supported as there is a go-yaml fork that uses
// `map[string]interface{}` instead.
func MergeMaps(
	src, tgt map[string]interface{}, itgt map[interface{}]interface{},
) {
	for sk, sv := range src {
		// 判断源key是否在目标字典存在，不存在则添加（判断不区分大小写）
		tk := KeyExists(sk, tgt)
		if tk == "" {
			tgt[sk] = sv
			if itgt != nil {
				itgt[sk] = sv
			}
			continue
		}

		// 如果源key在目标字典存在，但是大小写不匹配，需要添加到目标字典
		tv, ok := tgt[tk]
		if !ok {
			tgt[sk] = sv
			if itgt != nil {
				itgt[sk] = sv
			}
			continue
		}

		switch ttv := tv.(type) {
		case map[interface{}]interface{}:
			// 格式转换
			tsv, ok := sv.(map[interface{}]interface{})
			if !ok {
				continue
			}

			ssv := CastToMapStringInterface(tsv)
			stv := CastToMapStringInterface(ttv)
			MergeMaps(ssv, stv, ttv)
		case map[string]interface{}:
			tsv, ok := sv.(map[string]interface{})
			if !ok {
				continue
			}
			MergeMaps(tsv, ttv, nil)
		default:
			tgt[tk] = sv
			if itgt != nil {
				itgt[tk] = sv
			}
		}
	}
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

// CastToMapStringInterface 将字典中的key转换为字符串格式.
func CastToMapStringInterface(
	src map[interface{}]interface{},
) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}

// CastMapStringSliceToMapInterface 将字典的value值从 []string 转换为 interface{}.
func CastMapStringSliceToMapInterface(src map[string][]string) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[k] = v
	}
	return tgt
}

// CastMapStringToMapInterface 将字典的value值从 string 转换为 interface{}.
func CastMapStringToMapInterface(src map[string]string) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[k] = v
	}
	return tgt
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

// SetDeepMapValue 根据key设置DeepMap的值
func SetDeepMapValue(m map[string]interface{}, key string, value interface{}, keyDelim string) {
	key = strings.ToLower(key)
	value = ToCaseInsensitiveValue(value)

	path := strings.Split(key, keyDelim)
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := DeepSearch(m, path[0:len(path)-1])

	// set innermost value
	deepestMap[lastKey] = value
}

// SearchIndexableWithPathPrefixes recursively searches for a value for path in source map/slice.
//
// While SearchMap() considers each path element as a single map key or slice index, this
// function searches for, and prioritizes, merged path elements.
// e.g., if in the source, "foo" is defined with a sub-key "bar", and "foo.bar"
// is also defined, this latter value is returned for path ["foo", "bar"].
//
// This should be useful only at config level (other maps may not contain dots
// in their keys).
//
// Note: This assumes that the path entries and map keys are lower cased.
// 支持根据字典的key和数组的索引进行搜索
func SearchIndexableWithPathPrefixes(source interface{}, path []string, keyDelim string) interface{} {
	if len(path) == 0 {
		return source
	}

	// search for path prefixes, starting from the longest one
	for i := len(path); i > 0; i-- {
		// 获得当前路径和子路程
		prefixKey := strings.ToLower(strings.Join(path[0:i], keyDelim))

		var val interface{}
		switch sourceIndexable := source.(type) {
		case []interface{}:
			val = searchSliceWithPathPrefixes(sourceIndexable, prefixKey, i, path, keyDelim)
		case map[string]interface{}:
			val = searchMapWithPathPrefixes(sourceIndexable, prefixKey, i, path, keyDelim)
		}
		if val != nil {
			return val
		}
	}

	// not found
	return nil
}

// searchSliceWithPathPrefixes searches for a value for path in sourceSlice
//
// This function is part of the searchIndexableWithPathPrefixes recurring search and
// should not be called directly from functions other than searchIndexableWithPathPrefixes.
func searchSliceWithPathPrefixes(
	sourceSlice []interface{},
	prefixKey string,
	pathIndex int,
	path []string,
	keyDelim string,
) interface{} {
	// if the prefixKey is not a number or it is out of bounds of the slice
	index, err := strconv.Atoi(prefixKey)
	if err != nil || len(sourceSlice) <= index {
		return nil
	}

	next := sourceSlice[index]

	// Fast path
	if pathIndex == len(path) {
		return next
	}

	switch n := next.(type) {
	case map[interface{}]interface{}:
		return SearchIndexableWithPathPrefixes(cast.ToStringMap(n), path[pathIndex:], keyDelim)
	case map[string]interface{}, []interface{}:
		return SearchIndexableWithPathPrefixes(n, path[pathIndex:], keyDelim)
	default:
		// got a value but nested key expected, do nothing and look for next prefix
	}

	// not found
	return nil
}

// searchMapWithPathPrefixes searches for a value for path in sourceMap
//
// This function is part of the searchIndexableWithPathPrefixes recurring search and
// should not be called directly from functions other than searchIndexableWithPathPrefixes.
func searchMapWithPathPrefixes(
	sourceMap map[string]interface{},
	prefixKey string,
	pathIndex int,
	path []string,
	keyDelim string,
) interface{} {
	next, ok := sourceMap[prefixKey]
	if !ok {
		return nil
	}

	// Fast path
	if pathIndex == len(path) {
		return next
	}

	// Nested case
	switch n := next.(type) {
	case map[interface{}]interface{}:
		return SearchIndexableWithPathPrefixes(cast.ToStringMap(n), path[pathIndex:], keyDelim)
	case map[string]interface{}, []interface{}:
		return SearchIndexableWithPathPrefixes(n, path[pathIndex:], keyDelim)
	default:
		// got a value but nested key expected, do nothing and look for next prefix
	}

	// not found
	return nil
}

// IsPathShadowedInDeepMap makes sure the given path is not shadowed somewhere
// on its path in the map.
// e.g., if "foo.bar" has a value in the given map, it “shadows”
//       "foo.bar.baz" in a lower-priority map
// 判断子路径path是否覆盖到m中的一条路径，返回被覆盖的路径
func IsPathShadowedInDeepMap(path []string, m map[string]interface{}, keyDelim string) string {
	var parentVal interface{}
	for i := 1; i < len(path); i++ {
		parentVal = SearchMap(m, path[0:i])
		if parentVal == nil {
			// not found, no need to add more path elements
			return ""
		}
		switch parentVal.(type) {
		case map[interface{}]interface{}:
			continue
		case map[string]interface{}:
			continue
		default:
			// parentVal is a regular value which shadows "path"
			return strings.Join(path[0:i], keyDelim)
		}
	}
	return ""
}

// IsPathShadowedInFlatMap makes sure the given path is not shadowed somewhere
// in a sub-path of the map.
// e.g., if "foo.bar" has a value in the given map, it “shadows”
//       "foo.bar.baz" in a lower-priority map
// 判断子路径是否覆盖到m中的一条路径，返回被覆盖的路径
func IsPathShadowedInFlatMap(path []string, m map[string]interface{}, keyDelim string) string {
	// scan paths
	var parentKey string
	for i := 1; i < len(path); i++ {
		parentKey = strings.Join(path[0:i], keyDelim)
		if _, ok := m[parentKey]; ok {
			return parentKey
		}
	}
	return ""
}
