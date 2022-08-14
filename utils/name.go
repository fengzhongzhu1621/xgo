package utils

import (
	"fmt"
	"strings"
)

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
