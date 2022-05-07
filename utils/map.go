package utils

import (
	"fmt"
	"log"
	"reflect"
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

// DeepSearch THIS CODE IS COPIED HERE: IT SHOULD NOT BE MODIFIED
// AT SOME POINT IT WILL BE MOVED TO A COMMON PLACE
// deepSearch scans deep maps, following the key indexes listed in the
// sequence "path".
// The last value is expected to be another map, and is returned.
//
// In case intermediate keys do not exist, or map to a non-map value,
// a new map is created and inserted, and the search continues from there:
// the initial map "m" may be modified!
// 根据path构造字典，并返回指定路径上的最后一个字典，与 python 中的 setdefault 方法类似.
func DeepSearch(m map[string]interface{}, paths []string) map[string]interface{} {
	// 遍历 path 数组
	for _, k := range paths {
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

// mergeFlatMap merges the given maps, excluding values of the second map
// shadowed by values from the first map.
// MergeFlatMap 合并字典的key， value不可并；如果key重复，不合并到源字典.
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
			log.Println("", "tk", "\"\"", fmt.Sprintf("tgt[%s]", sk), sv)
			tgt[sk] = sv
			if itgt != nil {
				itgt[sk] = sv
			}
			continue
		}

		// 如果源key在目标字典存在，但是大小写不匹配，需要添加到目标字典
		tv, ok := tgt[tk]
		if !ok {
			log.Println("", fmt.Sprintf("ok[%s]", tk), false, fmt.Sprintf("tgt[%s]", sk), sv)
			tgt[sk] = sv
			if itgt != nil {
				itgt[sk] = sv
			}
			continue
		}

		svType := reflect.TypeOf(sv)
		tvType := reflect.TypeOf(tv)

		log.Println(
			"processing",
			"key", sk,
			"st", svType,
			"tt", tvType,
			"sv", sv,
			"tv", tv,
		)

		switch ttv := tv.(type) {
		case map[interface{}]interface{}:
			log.Println("merging maps (must convert)")
			// 格式转换
			tsv, ok := sv.(map[interface{}]interface{})
			if !ok {
				log.Printf(
					"Could not cast sv to map[interface{}]interface{}; key=%s, st=%v, tt=%v, sv=%v, tv=%v",
					sk, svType, tvType, sv, tv)
				continue
			}

			ssv := CastToMapStringInterface(tsv)
			stv := CastToMapStringInterface(ttv)
			MergeMaps(ssv, stv, ttv)
		case map[string]interface{}:
			log.Println("merging maps")
			tsv, ok := sv.(map[string]interface{})
			if !ok {
				log.Printf(
					"Could not cast sv to map[string]interface{}; key=%s, st=%v, tt=%v, sv=%v, tv=%v",
					sk, svType, tvType, sv, tv)
				continue
			}
			MergeMaps(tsv, ttv, nil)
		default:
			log.Println("setting value")
			tgt[tk] = sv
			if itgt != nil {
				itgt[tk] = sv
			}
		}
	}
}

// KeyExists 判断字段中是否存在指定key，key不区分大小写，返回源key
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

// CastToMapStringInterface 将字典中的key转换为字符串格式
func CastToMapStringInterface(
	src map[interface{}]interface{},
) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}

// CastMapStringSliceToMapInterface 将字典的value值从 []string 转换为 interface{}
func CastMapStringSliceToMapInterface(src map[string][]string) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[k] = v
	}
	return tgt
}

// CastMapStringToMapInterface 将字典的value值从 string 转换为 interface{}
func CastMapStringToMapInterface(src map[string]string) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[k] = v
	}
	return tgt
}
