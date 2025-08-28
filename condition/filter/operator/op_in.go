package operator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// InOp is in operator
type InOp OpType

// Name is in operator name
func (o InOp) Name() OpType {
	return In
}

// ValidateValue validate in operator's value
func (o InOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if opt == nil {
		return errors.New("validate option must be set")
	}

	err := validator.ValidateSliceOfBasicType(v, opt.MaxInLimit)
	if err != nil {
		return fmt.Errorf("in operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the in operator's field and value to a mongo query condition.
func (o InOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBIN: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o InOp) Match(value1, value2 interface{}) (bool, error) {
	var itemType string

	switch value1.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, json.Number:
		itemType = "numeric"
	case string:
		itemType = "string"
	case bool:
		itemType = "bool"
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("value(%+v) is not of basic type", value1)
	}

	if value2 == nil {
		return false, errors.New("rule value is nil")
	}

	switch reflect.TypeOf(value2).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return false, fmt.Errorf("rule value(%+v) is not of array type", value2)
	}

	v := reflect.ValueOf(value2)
	length := v.Len()
	if length == 0 {
		return false, errors.New("value is empty")
	}

	// 遍历数组，比较元素的值
	for i := 0; i < length; i++ {
		item := v.Index(i).Interface()

		switch itemType {
		case "numeric":
			val1, val2, err := parseNumericValues(value1, item)
			if err != nil {
				return false, err
			}
			if val1 == val2 {
				return true, nil
			}
		case "string":
			val, ok := item.(string)
			if !ok {
				return false, fmt.Errorf(
					"array ele index(%d) type(%T) not matches input type(%s)",
					i,
					item,
					itemType,
				)
			}
			if val == value1 {
				return true, nil
			}
		case "bool":
			val, ok := item.(bool)
			if !ok {
				return false, fmt.Errorf(
					"array ele index(%d) type(%T) not matches input type(%s)",
					i,
					item,
					itemType,
				)
			}
			if val == value1 {
				return true, nil
			}
		}
	}

	return false, nil
}
