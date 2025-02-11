package operator

import (
	"errors"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// LessOrEqualOp is less than or equal operator
type LessOrEqualOp OpType

// Name is less than or equal operator name
func (o LessOrEqualOp) Name() OpType {
	return LessOrEqual
}

// ValidateValue validate less than or equal operator value
func (o LessOrEqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !validator.IsNumeric(v) {
		return errors.New("invalid lte operator's value, should be a numeric value")
	}
	return nil
}

// ToMgo convert the less than or equal operator's field and value to a mongo query condition.
func (o LessOrEqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBLTE: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o LessOrEqualOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseNumericValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1 <= val2, nil
}
