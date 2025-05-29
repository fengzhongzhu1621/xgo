package reflectutils

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

// IsExported 判断结构体的字段名称首字符是否大写
func IsExported(field reflect.StructField) bool {
	// 获得字符串首字符
	r, _ := utf8.DecodeRuneInString(field.Name)
	// 判断是否是大写字符
	return r >= 'A' && r <= 'Z'
}

func IsExportedField(field *reflect.StructField) bool {
	// PkgPath 为空表示字段是导出的
	if field.PkgPath == "" {
		// 检查字段名的首字母是否为大写
		r, _ := utf8.DecodeRuneInString(field.Name)
		return unicode.IsUpper(r)
	}

	// PkgPath 不为空表示字段是未导出的
	return false
}

// IsExportedComponent 判断一个结构体字段是否为导出（exported）组件
func IsExportedComponent(field *reflect.StructField) bool {
	// field.PkgPath 表示字段所属的包路径
	pkgPath := field.PkgPath

	// 如果 len(pkgPath) > 0，意味着字段不是导出的（unexported）。
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

	// 检查字段名的首字母
	c := field.Name[0]

	// 如果字段名的首字母是大写字母（A-Z），则认为是导出的（exported）
	if 'a' <= c && c <= 'z' || c == '_' {
		return false
	}

	return true
}

// HasExportedFields 检查整个结构体是否包含至少一个导出字段
func HasExportedFields(t reflect.Type) bool {
	// 遍历所有字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 如果字段是导出的
		if field.PkgPath == "" {
			// 检查字段名的首字母是否为大写
			r, _ := utf8.DecodeRuneInString(field.Name)
			if unicode.IsUpper(r) {
				return true
			}
		}
		// 如果字段是嵌入的结构体，递归检查
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if HasExportedFields(field.Type) {
				return true
			}
		}
	}
	return false
}
