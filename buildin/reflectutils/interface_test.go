package reflectutils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInterface(t *testing.T) {
	var x float64 = 3.4
	t1 := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	fmt.Println("Type:", t1)             // Type: float64
	fmt.Println("Value:", v.Interface()) // Value: 3.4
}
