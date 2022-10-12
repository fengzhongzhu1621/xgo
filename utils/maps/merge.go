package maps

import (
	"strings"

	"github.com/fengzhongzhu1621/xgo/cast"
)

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
