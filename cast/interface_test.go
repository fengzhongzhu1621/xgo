package cast

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/duke-git/lancet/v2/convertor"
)

// Converts reflect value to its interface type.
// func ToInterface(v reflect.Value) (value interface{}, ok bool)
func TestToInterface(t *testing.T) {
	val := reflect.ValueOf("abc")
	iVal, ok := convertor.ToInterface(val)

	fmt.Printf("%T\n", iVal)
	fmt.Printf("%v\n", iVal)
	fmt.Println(ok)

	// Output:
	// string
	// abc
	// true
}
