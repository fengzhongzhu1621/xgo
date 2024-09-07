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

func IsExportedComponent(field *reflect.StructField) bool {
	pkgPath := field.PkgPath
	if len(pkgPath) > 0 {
		// 结构体字段不是一个struct类型
		// 普通字段，例如 Manager.title
		return false
	}
	// 结构体字段是一个struct类型，例如 Manager中的User
	// type User struct {
	//     Id   int
	//	   Name string
	//	   Age  int
	// }
	//
	// type Manager struct {
	//     User
	//	   title string
	// }
	c := field.Name[0]
	if 'a' <= c && c <= 'z' || c == '_' {
		return false
	}
	// 首字母大写
	return true
}
