package reflectutils

import (
	"reflect"
)

// Indirect 获得反射值对象
func Indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

// IndirectType 获得反射值类型，并判断其是否是指针类型
func IndirectType(reflectType reflect.Type) (_ reflect.Type, isPtr bool) {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		// 返回指针指向的值的类型
		reflectType = reflectType.Elem()
		isPtr = true
	}
	return reflectType, isPtr
}
