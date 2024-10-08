package cast

import (
	"errors"
	"fmt"
	"reflect"
)

// ToSlice casts an interface to a []interface{} type.
func ToSlice(i interface{}) []interface{} {
	v, _ := ToSliceE(i)
	return v
}

// ToSliceE casts an interface to a []interface{} type.
func ToSliceE(i interface{}) ([]interface{}, error) {
	var s []interface{}

	switch v := i.(type) {
	case []interface{}:
		// 输入类型是数组，则采用数组拼接的方式
		return append(s, v...), nil
	case []map[string]interface{}:
		// 输入类型是字典，则获得所有 value 组成的数组
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to []interface{}", i, i)
	}
}

var ErrNotArray = errors.New("only support array")

// ToSlice conv an array-interface to []interface{}
// will error if the type is not slice
func ToSlice2(array interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(array)
	if v.Kind() != reflect.Slice {
		return nil, ErrNotArray
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret, nil
}

func TransSlice2Interface(old interface{}) ([]interface{}, error) {
	switch trans := old.(type) {
	case []string:
		new := make([]interface{}, len(trans))
		for i, v := range trans {
			new[i] = v
		}
		return new, nil
	default:
		return nil, errors.New("illegal type")
	}
}

func TransSlice2Interface2(slice []string) ([]interface{}, error) {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice, nil
}
