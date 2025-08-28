package validator

import (
	"fmt"
	"reflect"
)

// ValidateBoolType validate if the value is a bool type
func ValidateBoolType(value interface{}) error {
	if reflect.TypeOf(value).Kind() != reflect.Bool {
		return fmt.Errorf("value(%+v) is not of bool type", value)
	}
	return nil
}
