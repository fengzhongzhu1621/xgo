package validator

import (
	"fmt"
	"reflect"
)

// MustBeFunction 判断是否为函数类型
func MustBeFunction(function any) {
	v := reflect.ValueOf(function)
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("Invalid function type, value of type %T", function))
	}
}
