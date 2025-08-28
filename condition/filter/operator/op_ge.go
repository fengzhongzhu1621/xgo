package operator

import (
	"errors"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// GreaterOrEqualOp is greater than or equal operator
type GreaterOrEqualOp OpType

// Name is greater than or equal operator name
func (o GreaterOrEqualOp) Name() OpType {
	return GreaterOrEqual
}

// ValidateValue validate greater than or equal operator value
func (o GreaterOrEqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !validator.IsNumeric(v) {
		return errors.New("invalid gte operator's value, should be a numeric value")
	}
	return nil
}

// ToMgo convert the greater than or equal operator's field and value to a mongo query condition.
func (o GreaterOrEqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBGTE: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o GreaterOrEqualOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseNumericValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1 >= val2, nil
}
