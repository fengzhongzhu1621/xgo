package operator

import (
	"errors"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// NotEqualOp is not equal operator type
type NotEqualOp OpType

// Name is not equal operator name
func (ne NotEqualOp) Name() OpType {
	return NotEqual
}

// ValidateValue validate not equal operator's value
func (ne NotEqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if !validator.IsBasicValue(v) {
		return fmt.Errorf("invalid ne value(%+v)", v)
	}
	return nil
}

// ToMgo convert the not equal operator's field and value to a mongo query condition.
func (ne NotEqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBNE: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (ne NotEqualOp) Match(value1, value2 interface{}) (bool, error) {
	matched, err := GetOperator(Equal).Match(value1, value2)
	if err != nil {
		return false, err
	}
	return !matched, nil
}
