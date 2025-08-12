// Package masking 敏感信息脱敏拦截器
package masking

import (
	"reflect"
	"time"
)

// DeepCheck 递归处理待脱敏数据
func DeepCheck(src interface{}) {
	if src == nil {
		return
	}

	original := reflect.ValueOf(src)

	checkRecursive(original)
}

// checkRecursive 调用masking插件导出的脱敏方法
func checkRecursive(original reflect.Value) {
	// 判断值的类型的枚举
	switch original.Kind() {
	case reflect.Ptr:
		originalValue := original.Elem()

		// 如果 reflect.Value 的 IsValid() 方法返回 false，
		// 那么它就是一个无效的反射对象，调用它的任何方法都会 panic，除了 String 方法。
		if !originalValue.IsValid() {
			return
		}
		checkRecursive(originalValue)

	case reflect.Interface:
		if original.IsNil() {
			return
		}
		originalValue := original.Elem()
		if v, ok := originalValue.Interface().(Masking); ok {
			v.Masking()
		}
		checkRecursive(originalValue)

	case reflect.Struct:
		// 将 reflect.Value 转换回普通的 interface{} 类型，以便作为函数的返回值。
		_, ok := original.Interface().(time.Time)
		if ok {
			return
		}
		for i := 0; i < original.NumField(); i++ {
			if original.Type().Field(i).PkgPath != "" {
				continue
			}
			if v, ok := original.Field(i).Interface().(Masking); ok {
				v.Masking()
			}
			checkRecursive(original.Field(i))
		}

	case reflect.Slice:
		if original.IsNil() {
			return
		}
		for i := 0; i < original.Len(); i++ {
			if v, ok := original.Index(i).Interface().(Masking); ok {
				v.Masking()
			}
			checkRecursive(original.Index(i))
		}

	case reflect.Map:
		if original.IsNil() {
			return
		}
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			checkRecursive(originalValue)
			if v, ok := key.Interface().(Masking); ok {
				v.Masking()
			}
		}
	}
}
