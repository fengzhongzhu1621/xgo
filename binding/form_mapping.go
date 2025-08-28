// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"fmt"
	"reflect"

	reflect_utils "github.com/fengzhongzhu1621/xgo/buildin/reflectutils"
	string_utils "github.com/fengzhongzhu1621/xgo/str/stringutils"
)

type setOptions struct {
	isDefaultExists bool // 表单字段是否存在默认值
	defaultValue    string
}

// setter tries to set value on a walking by fields of a struct.
type setter interface {
	TrySet(
		value reflect.Value,
		field reflect.StructField,
		key string,
		opt setOptions,
	) (isSetted bool, err error)
}

// 定义表单对象.
type formSource map[string][]string

// 验证表单有setter接口.
var _ setter = formSource(nil)

var emptyField = reflect.StructField{}

// TrySet tries to set a value by request's form source (like map[string][]string).
func (form formSource) TrySet(
	value reflect.Value,
	field reflect.StructField,
	tagValue string,
	opt setOptions,
) (isSetted bool, err error) {
	return setByForm(value, field, form, tagValue, opt)
}

func mappingByPtr(ptr interface{}, setter setter, tag string) error {
	_, err := mapping(reflect.ValueOf(ptr), emptyField, setter, tag)
	return err
}

func mapFormByTag(ptr interface{}, form map[string][]string, tag string) error {
	return mappingByPtr(ptr, formSource(form), tag)
}

func mapUri(ptr interface{}, m map[string][]string) error {
	return mapFormByTag(ptr, m, "uri")
}

func mapForm(ptr interface{}, form map[string][]string) error {
	return mapFormByTag(ptr, form, "form")
}

func mapping(
	value reflect.Value,
	field reflect.StructField,
	setter setter,
	tag string,
) (bool, error) {
	if field.Tag.Get(tag) == "-" { // just ignoring this field
		return false, nil
	}

	vKind := value.Kind()

	if vKind == reflect.Ptr {
		var isNew bool
		vPtr := value
		if value.IsNil() {
			isNew = true
			vPtr = reflect.New(value.Type().Elem())
		}
		isSetted, err := mapping(vPtr.Elem(), field, setter, tag)
		if err != nil {
			return false, err
		}
		if isNew && isSetted {
			value.Set(vPtr)
		}
		return isSetted, nil
	}

	if vKind != reflect.Struct || !field.Anonymous {
		ok, err := tryToSetValue(value, field, setter, tag)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	if vKind == reflect.Struct {
		tValue := value.Type()

		var isSetted bool
		for i := 0; i < value.NumField(); i++ {
			sf := tValue.Field(i)
			if sf.PkgPath != "" && !sf.Anonymous { // unexported
				continue
			}
			ok, err := mapping(value.Field(i), tValue.Field(i), setter, tag)
			if err != nil {
				return false, err
			}
			isSetted = isSetted || ok
		}
		return isSetted, nil
	}
	return false, nil
}

func tryToSetValue(
	value reflect.Value,
	field reflect.StructField,
	setter setter,
	tag string,
) (bool, error) {
	var tagValue string
	var setOpt setOptions

	// 获得表单tag，例如 `form:"check_in" binding:"required" time_format:"2006-01-02"`，tagValue值为check_in
	tagValue = field.Tag.Get(tag)
	tagValue, opts := string_utils.Head(tagValue, ",")

	if tagValue == "" { // default value is FieldName
		tagValue = field.Name
	}
	if tagValue == "" { // when field is "emptyField" variable
		return false, nil
	}

	var opt string
	for len(opts) > 0 {
		opt, opts = string_utils.Head(opts, ",")

		if k, v := string_utils.Head(opt, "="); k == "default" {
			setOpt.isDefaultExists = true
			setOpt.defaultValue = v
		}
	}

	return setter.TrySet(value, field, tagValue, setOpt)
}

// 根据表单的值填充对象.
func setByForm(
	value reflect.Value,
	field reflect.StructField,
	form map[string][]string,
	tagValue string,
	opt setOptions,
) (isSetted bool, err error) {
	// 获得表单的值
	vs, ok := form[tagValue]
	if !ok && !opt.isDefaultExists {
		return false, nil
	}

	switch value.Kind() {
	case reflect.Slice:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		return true, reflect_utils.SetSlice(vs, value, field)
	case reflect.Array:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		if len(vs) != value.Len() {
			return false, fmt.Errorf("%q is not valid value for %s", vs, value.Type().String())
		}
		return true, reflect_utils.SetArray(vs, value, field)
	default:
		var val string
		if !ok {
			val = opt.defaultValue
		}

		if len(vs) > 0 {
			val = vs[0]
		}
		return true, reflect_utils.SetWithProperType(val, value, field)
	}
}
