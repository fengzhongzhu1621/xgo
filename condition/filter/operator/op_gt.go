package operator

import (
	"errors"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// GreaterOp is greater than operator
type GreaterOp OpType

// Name is greater than operator name
func (o GreaterOp) Name() OpType {
	return Greater
}

// ValidateValue validate greater than operator value
func (o GreaterOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !validator.IsNumeric(v) {
		return errors.New("invalid gt operator's value, should be a numeric value")
	}
	return nil
}

// ToMgo convert the greater than operator's field and value to a mongo query condition.
func (o GreaterOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBGT: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o GreaterOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseNumericValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1 > val2, nil
}
