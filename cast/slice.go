package cast

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

// ToSlice2 conv an array-interface to []interface{}
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

// ToStringSlice casts an interface to a []string type.
func ToStringSlice(i interface{}) []string {
	v, _ := ToStringSliceE(i)
	return v
}

// ToStringSliceE casts an interface to a []string type.
func ToStringSliceE(i interface{}) ([]string, error) {
	var a []string

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []string:
		return v, nil
	case []int8:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int32:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int64:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []float32:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []float64:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case string:
		return strings.Fields(v), nil
	case []error:
		for _, err := range i.([]error) {
			a = append(a, err.Error())
		}
		return a, nil
	case interface{}:
		str, err := ToStringE(v)
		if err != nil {
			return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
		}
		return []string{str}, nil
	default:
		return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
	}
}

// //////////////////////////////////////////////////////////////////////////////////////
// StringToInt64Slice 根据分隔符将字符串转换为整型数组
// 1,2,3 -> []int64{1, 2, 3}
func StringToInt64Slice(s, sep string) ([]int64, error) {
	if s == "" {
		return []int64{}, nil
	}
	parts := strings.Split(s, sep)

	int64Slice := make([]int64, 0, len(parts))
	for _, d := range parts {
		i, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			return nil, err
		}
		int64Slice = append(int64Slice, i)
	}
	return int64Slice, nil
}

// StringToIntSlice 根据分隔符将字符串转换为整型数组
// 1,2,3 -> []int64{1, 2, 3}
func StringToIntSlice(s, sep string) ([]int, error) {
	if s == "" {
		return []int{}, nil
	}
	parts := strings.Split(s, sep)

	intSlice := make([]int, 0, len(parts))
	for _, d := range parts {
		i, err := strconv.Atoi(d)
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, i)
	}
	return intSlice, nil
}
