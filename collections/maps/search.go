package maps

import (
	"strconv"
	"strings"

	"github.com/fengzhongzhu1621/xgo/cast"
)

// SearchMap recursively searches for a value for path in source map.
// Returns nil if not found.
// Note: This assumes that the path entries and map keys are lower cased.
// 在m中搜索路径path.
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
		switch next := next.(type) {
		case map[interface{}]interface{}:
			return SearchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return SearchMap(next, path[1:])
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
// 支持根据字典的key和数组的索引进行搜索.
func SearchIndexableWithPathPrefixes(
	source interface{},
	path []string,
	keyDelim string,
) interface{} {
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
//
//	"foo.bar.baz" in a lower-priority map
//
// 判断子路径path是否覆盖到m中的一条路径，返回被覆盖的路径.
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
//
//	"foo.bar.baz" in a lower-priority map
//
// 判断子路径是否覆盖到m中的一条路径，返回被覆盖的路径.
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
