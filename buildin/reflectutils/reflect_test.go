package reflectutils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypeOf(t *testing.T) {
	// reflect.TypeOf(map[string]interface{}(nil))返回一个表示map[string]interface{}类型的reflect.Type值。
	tp := reflect.TypeOf(map[string]interface{}(nil))
	fmt.Println(tp) // map[string]interface {}
}
