package buildin

import (
	"fmt"
	"strings"
)

// FullyQualifiedStructName name returns object name in format [package].[type name].
// It ignores if the value is a pointer or not.
func FullyQualifiedStructName(v interface{}) string {
	// 获得值的类型
	s := fmt.Sprintf("%T", v)
	s = strings.TrimLeft(s, "*")

	return s
}

// StructName returns a normalized name of the passed structure.
// 获得结构体的名称
func StructName(v interface{}) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}

	// %T      相应值的类型的Go语法表示       Printf("%T", people)   main.Human
	s := fmt.Sprintf("%T", v)
	// trim the pointer marker, if any
	return strings.TrimLeft(s, "*")
}

// StructName name returns struct name in format [type name].
// It ignores if the value is a pointer or not.
func RawStructName(v interface{}) string {
	segments := strings.Split(fmt.Sprintf("%T", v), ".")

	return segments[len(segments)-1]
}

type namedStruct interface {
	Name() string
}

// NamedStruct returns the name from a message implementing the following interface:
// 		type namedStruct interface {
// 			Name() string
// 		}
// It ignores if the value is a pointer or not.
func NamedStruct(fallback func(v interface{}) string) func(v interface{}) string {
	return func(v interface{}) string {
		if v, ok := v.(namedStruct); ok {
			return v.Name()
		}

		return fallback(v)
	}
}
