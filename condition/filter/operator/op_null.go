package operator

import (
	"errors"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
)

// IsNullOp is null operator
type IsNullOp OpType

// Name is null operator name
func (o IsNullOp) Name() OpType {
	return IsNull
}

// ValidateValue validate null operator's value
func (o IsNullOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the null operator's field and value to a mongo query condition.
func (o IsNullOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is null")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBEQ: nil,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o IsNullOp) Match(value1, value2 interface{}) (bool, error) {
	return value1 == nil, nil
}

// IsNotNullOp is not null operator
type IsNotNullOp OpType

// Name is not null operator name
func (o IsNotNullOp) Name() OpType {
	return IsNotNull
}

// ValidateValue validate is not null operator's value
func (o IsNotNullOp) ValidateValue(v interface{}, opt *ExprOption) error {
	return nil
}

// ToMgo convert the is not null operator's field and value to a mongo query condition.
func (o IsNotNullOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is null")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBNE: nil,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o IsNotNullOp) Match(value1, value2 interface{}) (bool, error) {
	return value1 != nil, nil
}
