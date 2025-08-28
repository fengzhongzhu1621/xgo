package operator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
)

// IsEmptyOp is empty operator
type IsEmptyOp OpType

// Name is empty operator name
func (o IsEmptyOp) Name() OpType {
	return IsEmpty
}

// ValidateValue validate empty operator's value
func (o IsEmptyOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the empty operator's field and value to a mongo query condition.
func (o IsEmptyOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBSize: 0,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o IsEmptyOp) Match(value1, value2 interface{}) (bool, error) {
	if value1 == nil {
		return false, errors.New("input value is nil")
	}

	switch reflect.TypeOf(value1).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return false, fmt.Errorf("rule value(%+v) is not of array type", value1)
	}

	return reflect.ValueOf(value1).Len() == 0, nil
}

// IsNotEmptyOp is not empty operator
type IsNotEmptyOp OpType

// Name is not empty operator name
func (o IsNotEmptyOp) Name() OpType {
	return IsNotEmpty
}

// ValidateValue validate is not empty operator's value
func (o IsNotEmptyOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the is not empty operator's field and value to a mongo query condition.
func (o IsNotEmptyOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBSize: map[string]interface{}{DBGT: 0},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o IsNotEmptyOp) Match(value1, value2 interface{}) (bool, error) {
	if value1 == nil {
		return false, errors.New("input value is nil")
	}

	switch reflect.TypeOf(value1).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return false, fmt.Errorf("rule value(%+v) is not of array type", value1)
	}

	return reflect.ValueOf(value1).Len() > 0, nil
}
