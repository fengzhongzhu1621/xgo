package reflectutils

import (
	"reflect"
	"unicode/utf8"
)

// 判断结构体的字段名称首字符是否大写
func IsExported(field reflect.StructField) bool {
	// 获得字符串首字符
	r, _ := utf8.DecodeRuneInString(field.Name)
	// 判断是否是大写字符
	return r >= 'A' && r <= 'Z'
}
