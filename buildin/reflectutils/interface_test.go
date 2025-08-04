package reflectutils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInterface(t *testing.T) {
	x := 3.4
	t1 := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	fmt.Println("Type:", t1)             // Type: float64
	fmt.Println("Value:", v.Interface()) // Value: 3.4
}

// 反射允许在运行时对接口变量进行类型断言
func TestInterfaceAny(t *testing.T) {
	var i interface{} = "hello"
	v := reflect.ValueOf(i)

	if s, ok := v.Interface().(string); ok {
		fmt.Println(s) // hello
	}
}
