package reflectutils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSetFloat(t *testing.T) {
	x := 3.14

	// Elem(): 它允许你访问指针指向的值或接口包含的值。
	// 接受一个 reflect.Value 类型的参数，并返回一个新的 reflect.Value，表示指针指向的值或接口包含的值。
	// 如果传入的值不是指针或接口类型，Elem() 函数将引发 panic。
	v := reflect.ValueOf(&x).Elem() // 返回给定值的值
	v.SetFloat(3.1415926)
	fmt.Println("x: ", v) // x:  3.1415926
}
