package operator

import (
	"errors"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// LessOp is less than operator
type LessOp OpType

// Name is less than operator name
func (o LessOp) Name() OpType {
	return Less
}

// ValidateValue validate less than operator value
func (o LessOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !validator.IsNumeric(v) {
		return fmt.Errorf("invalid lt operator's value, should be a numeric value")
	}
	return nil
}

// ToMgo convert the less than  operator's field and value to a mongo query condition.
func (o LessOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBLT: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o LessOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseNumericValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1 < val2, nil
}
