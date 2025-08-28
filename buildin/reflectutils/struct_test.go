package reflectutils

import (
	"fmt"
	"reflect"
	"testing"
)

type MyMethods struct{}

func (m *MyMethods) MyMethod() string {
	return "Hello, World!"
}

type Manager struct {
	User
	title string
}

func TestStructHasExportedFields(t *testing.T) {
	managerType := reflect.TypeOf(Manager{})
	hasExported := HasExportedFields(managerType)
	fmt.Println(hasExported) // 输出: true
}

// 访问结构体字段
func TestStructFieldByName(t *testing.T) {
	type MyStruct struct {
		privateField int
	}

	s := MyStruct{privateField: 1}
	v := reflect.ValueOf(s)

	field := v.FieldByName("privateField")
	fmt.Println("Private Field:", field.Int()) // Private Field: 1
}

func TestStructMethodByName(t *testing.T) {
	obj := &MyMethods{}
	v := reflect.ValueOf(obj)
	method := v.MethodByName("MyMethod")
	result := method.Call(nil)

	fmt.Println("Method Result:", result[0].Interface()) // Method Result: Hello, World!
}
