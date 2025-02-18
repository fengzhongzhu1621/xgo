package operator

import (
	"errors"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
)

// ExistOp is 'exist' operator
type ExistOp OpType

// Name is 'exist' operator name
func (o ExistOp) Name() OpType {
	return Exist
}

// ValidateValue validate 'exist' operator's value
func (o ExistOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the 'exist' operator's field and value to a mongo query condition.
func (o ExistOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is null")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBExists: true,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o ExistOp) Match(value1, value2 interface{}) (bool, error) {
	return value1 == nil, nil
}

// NotExistOp is not exist operator
type NotExistOp OpType

// Name is not exist operator name
func (o NotExistOp) Name() OpType {
	return NotExist
}

// ValidateValue validate is not exist operator's value
func (o NotExistOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the is not exist operator's field and value to a mongo query condition.
func (o NotExistOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is null")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBExists: false,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotExistOp) Match(value1, value2 interface{}) (bool, error) {
	return value1 != nil, nil
}
