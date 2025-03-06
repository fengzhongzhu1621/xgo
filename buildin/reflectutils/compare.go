package reflectutils

import "reflect"

func IsEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
