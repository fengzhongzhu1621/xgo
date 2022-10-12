package maps

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/fengzhongzhu1621/xgo/utils/reflectutils"
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

func hasMergeableFields(dst reflect.Value) (exported bool) {
	// 遍历结构体所有的字段
	for i, n := 0, dst.NumField(); i < n; i++ {
		// 获得字段
		field := dst.Type().Field(i)
		if field.Anonymous && dst.Field(i).Kind() == reflect.Struct {
			exported = exported || hasMergeableFields(dst.Field(i))
		} else if reflectutils.IsExportedComponent(&field) {
			exported = exported || len(field.PkgPath) == 0
		}
	}
	return
}

// Traverses recursively both values, assigning src's fields values to dst.
// The map argument tracks comparisons that have already been seen, which allows
// short circuiting on recursive types.
func deepMerge(dst, src reflect.Value, visited map[uintptr]*visit, depth int, config *Config) (err error) {
	overwrite := config.Overwrite
	typeCheck := config.TypeCheck
	overwriteWithEmptySrc := config.overwriteWithEmptyValue
	overwriteSliceWithEmptySrc := config.overwriteSliceWithEmptyValue
	sliceDeepCopy := config.sliceDeepCopy

	if !src.IsValid() {
		return
	}
	if dst.CanAddr() {
		addr := dst.UnsafeAddr()
		h := 17 * addr
		seen := visited[h]
		typ := dst.Type()
		for p := seen; p != nil; p = p.next {
			if p.ptr == addr && p.typ == typ {
				return nil
			}
		}
		// Remember, remember...
		visited[h] = &visit{typ, seen, addr}
	}

	if config.Transformers != nil && !reflectutils.IsReflectNil(dst) && dst.IsValid() {
		if fn := config.Transformers.Transformer(dst.Type()); fn != nil {
			err = fn(dst, src)
			return
		}
	}

	switch dst.Kind() {
	case reflect.Struct:
		if hasMergeableFields(dst) {
			for i, n := 0, dst.NumField(); i < n; i++ {
				if err = deepMerge(dst.Field(i), src.Field(i), visited, depth+1, config); err != nil {
					return
				}
			}
		} else {
			if dst.CanSet() && (reflectutils.IsReflectNil(dst) || overwrite) && (!reflectutils.IsEmptyValue(src) || overwriteWithEmptySrc) {
				dst.Set(src)
			}
		}
	case reflect.Map:
		if dst.IsNil() && !src.IsNil() {
			if dst.CanSet() {
				dst.Set(reflect.MakeMap(dst.Type()))
			} else {
				dst = src
				return
			}
		}

		if src.Kind() != reflect.Map {
			if overwrite {
				dst.Set(src)
			}
			return
		}

		for _, key := range src.MapKeys() {
			srcElement := src.MapIndex(key)
			if !srcElement.IsValid() {
				continue
			}
			dstElement := dst.MapIndex(key)
			switch srcElement.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Interface, reflect.Slice:
				if srcElement.IsNil() {
					if overwrite {
						dst.SetMapIndex(key, srcElement)
					}
					continue
				}
				fallthrough
			default:
				if !srcElement.CanInterface() {
					continue
				}
				switch reflect.TypeOf(srcElement.Interface()).Kind() {
				case reflect.Struct:
					fallthrough
				case reflect.Ptr:
					fallthrough
				case reflect.Map:
					srcMapElm := srcElement
					dstMapElm := dstElement
					if srcMapElm.CanInterface() {
						srcMapElm = reflect.ValueOf(srcMapElm.Interface())
						if dstMapElm.IsValid() {
							dstMapElm = reflect.ValueOf(dstMapElm.Interface())
						}
					}
					if err = deepMerge(dstMapElm, srcMapElm, visited, depth+1, config); err != nil {
						return
					}
				case reflect.Slice:
					srcSlice := reflect.ValueOf(srcElement.Interface())

					var dstSlice reflect.Value
					if !dstElement.IsValid() || dstElement.IsNil() {
						dstSlice = reflect.MakeSlice(srcSlice.Type(), 0, srcSlice.Len())
					} else {
						dstSlice = reflect.ValueOf(dstElement.Interface())
					}

					if (!reflectutils.IsEmptyValue(src) || overwriteWithEmptySrc || overwriteSliceWithEmptySrc) && (overwrite || reflectutils.IsEmptyValue(dst)) && !config.AppendSlice && !sliceDeepCopy {
						if typeCheck && srcSlice.Type() != dstSlice.Type() {
							return fmt.Errorf("cannot override two slices with different type (%s, %s)", srcSlice.Type(), dstSlice.Type())
						}
						dstSlice = srcSlice
					} else if config.AppendSlice {
						if srcSlice.Type() != dstSlice.Type() {
							return fmt.Errorf("cannot append two slices with different type (%s, %s)", srcSlice.Type(), dstSlice.Type())
						}
						dstSlice = reflect.AppendSlice(dstSlice, srcSlice)
					} else if sliceDeepCopy {
						i := 0
						for ; i < srcSlice.Len() && i < dstSlice.Len(); i++ {
							srcElement := srcSlice.Index(i)
							dstElement := dstSlice.Index(i)

							if srcElement.CanInterface() {
								srcElement = reflect.ValueOf(srcElement.Interface())
							}
							if dstElement.CanInterface() {
								dstElement = reflect.ValueOf(dstElement.Interface())
							}

							if err = deepMerge(dstElement, srcElement, visited, depth+1, config); err != nil {
								return
							}
						}

					}
					dst.SetMapIndex(key, dstSlice)
				}
			}
			if dstElement.IsValid() && !reflectutils.IsEmptyValue(dstElement) && (reflect.TypeOf(srcElement.Interface()).Kind() == reflect.Map || reflect.TypeOf(srcElement.Interface()).Kind() == reflect.Slice) {
				continue
			}

			if srcElement.IsValid() && ((srcElement.Kind() != reflect.Ptr && overwrite) || !dstElement.IsValid() || reflectutils.IsEmptyValue(dstElement)) {
				if dst.IsNil() {
					dst.Set(reflect.MakeMap(dst.Type()))
				}
				dst.SetMapIndex(key, srcElement)
			}
		}
	case reflect.Slice:
		if !dst.CanSet() {
			break
		}
		if (!reflectutils.IsEmptyValue(src) || overwriteWithEmptySrc || overwriteSliceWithEmptySrc) && (overwrite || reflectutils.IsEmptyValue(dst)) && !config.AppendSlice && !sliceDeepCopy {
			dst.Set(src)
		} else if config.AppendSlice {
			if src.Type() != dst.Type() {
				return fmt.Errorf("cannot append two slice with different type (%s, %s)", src.Type(), dst.Type())
			}
			dst.Set(reflect.AppendSlice(dst, src))
		} else if sliceDeepCopy {
			for i := 0; i < src.Len() && i < dst.Len(); i++ {
				srcElement := src.Index(i)
				dstElement := dst.Index(i)
				if srcElement.CanInterface() {
					srcElement = reflect.ValueOf(srcElement.Interface())
				}
				if dstElement.CanInterface() {
					dstElement = reflect.ValueOf(dstElement.Interface())
				}

				if err = deepMerge(dstElement, srcElement, visited, depth+1, config); err != nil {
					return
				}
			}
		}
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		if reflectutils.IsReflectNil(src) {
			if overwriteWithEmptySrc && dst.CanSet() && src.Type().AssignableTo(dst.Type()) {
				dst.Set(src)
			}
			break
		}

		if src.Kind() != reflect.Interface {
			if dst.IsNil() || (src.Kind() != reflect.Ptr && overwrite) {
				if dst.CanSet() && (overwrite || reflectutils.IsEmptyValue(dst)) {
					dst.Set(src)
				}
			} else if src.Kind() == reflect.Ptr {
				if err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, config); err != nil {
					return
				}
			} else if dst.Elem().Type() == src.Type() {
				if err = deepMerge(dst.Elem(), src, visited, depth+1, config); err != nil {
					return
				}
			} else {
				return ErrDifferentArgumentsTypes
			}
			break
		}

		if dst.IsNil() || overwrite {
			if dst.CanSet() && (overwrite || reflectutils.IsEmptyValue(dst)) {
				dst.Set(src)
			}
			break
		}

		if dst.Elem().Kind() == src.Elem().Kind() {
			if err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, config); err != nil {
				return
			}
			break
		}
	default:
		mustSet := (reflectutils.IsEmptyValue(dst) || overwrite) && (!reflectutils.IsEmptyValue(src) || overwriteWithEmptySrc)
		if mustSet {
			if dst.CanSet() {
				dst.Set(src)
			} else {
				dst = src
			}
		}
	}

	return
}

func merge(dst, src interface{}, opts ...func(*Config)) error {
	if dst != nil && reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return ErrNonPointerAgument
	}
	var (
		vDst, vSrc reflect.Value
		err        error
	)

	config := &Config{}

	for _, opt := range opts {
		opt(config)
	}

	if vDst, vSrc, err = resolveValues(dst, src); err != nil {
		return err
	}
	if vDst.Type() != vSrc.Type() {
		return ErrDifferentArgumentsTypes
	}
	return deepMerge(vDst, vSrc, make(map[uintptr]*visit), 0, config)
}

// Merge will fill any empty for value type attributes on the dst struct using corresponding
// src attributes if they themselves are not empty. dst and src must be valid same-type structs
// and dst must be a pointer to struct.
// It won't merge unexported (private) fields and will do recursively any exported field.
func Merge(dst, src interface{}, opts ...func(*Config)) error {
	return merge(dst, src, opts...)
}
