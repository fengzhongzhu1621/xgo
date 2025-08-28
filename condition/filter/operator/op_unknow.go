package operator

import "errors"

// UnknownOp is unknown operator
type UnknownOp OpType

// Name is equal operator
func (o UnknownOp) Name() OpType {
	return Unknown
}

// ValidateValue validate equal's value
func (o UnknownOp) ValidateValue(_ interface{}, _ *ExprOption) error {
	return errors.New("unknown operator")
}

// ToMgo convert this operator's field and value to a mongo query condition.
func (o UnknownOp) ToMgo(_ string, _ interface{}) (map[string]interface{}, error) {
	return nil, errors.New("unknown operator, can not gen mongo expression")
}

// Match checks if the first data matches the second data by this operator
func (o UnknownOp) Match(_, _ interface{}) (bool, error) {
	return false, errors.New("unknown operator, can not check if two value matches this operator")
}
