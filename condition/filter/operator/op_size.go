package operator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
)

// SizeOp size operator
type SizeOp OpType

// Name size operator name
func (o SizeOp) Name() OpType {
	return Size
}

// ValidateValue validate size operator's value
func (o SizeOp) ValidateValue(v interface{}, opt *ExprOption) error {
	intVal, err := cast.ToInt64E(v)
	if err != nil {
		return fmt.Errorf("invalid size operator's value, should be a numeric value, err: %v", err)
	}

	if intVal < 0 {
		return fmt.Errorf("invalid size operator's value, should not be negative")
	}
	return nil
}

// ToMgo convert the size operator's field and value to a mongo query condition.
func (o SizeOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBSize: value,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o SizeOp) Match(value1, value2 interface{}) (bool, error) {
	if value1 == nil {
		return false, errors.New("input value is nil")
	}

	switch reflect.TypeOf(value1).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return false, fmt.Errorf("rule value(%+v) is not of array type", value1)
	}

	intVal, err := cast.ToIntE(value2)
	if err != nil {
		return false, fmt.Errorf(
			"invalid size operator's value, should be a numeric value, err: %v",
			err,
		)
	}

	return reflect.ValueOf(value1).Len() == intVal, nil
}
