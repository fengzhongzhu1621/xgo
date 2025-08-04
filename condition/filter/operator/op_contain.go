package operator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// ContainsOp is contains operator
type ContainsOp OpType

// Name is contains operator name
func (o ContainsOp) Name() OpType {
	return Contains
}

// ValidateValue validate contains operator's value
func (o ContainsOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("contains operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the contains operator's field and value to a mongo query condition.
func (o ContainsOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBLIKE:    value,
			DBOPTIONS: "i",
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o ContainsOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.Contains(strings.ToLower(val1), strings.ToLower(val2)), nil
}

// ContainsSensitiveOp is contains sensitive operator
type ContainsSensitiveOp OpType

// Name is contains sensitive operator name
func (o ContainsSensitiveOp) Name() OpType {
	return ContainsSensitive
}

// ValidateValue validate contains sensitive operator's value
func (o ContainsSensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("contains sensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the contains sensitive operator's field and value to a mongo query condition.
func (o ContainsSensitiveOp) ToMgo(
	field string,
	value interface{},
) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBLIKE: value,
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o ContainsSensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.Contains(val1, val2), nil
}

// NotContainsOp is not contains operator
type NotContainsOp OpType

// Name is not contains operator name
func (o NotContainsOp) Name() OpType {
	return NotContains
}

// ValidateValue validate not contains operator's value
func (o NotContainsOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not contains operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not contains operator's field and value to a mongo query condition.
func (o NotContainsOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBNot: map[string]interface{}{
				DBLIKE: value,
			},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotContainsOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.Contains(val1, val2), nil
}

// NotContainsInsensitiveOp is not contains insensitive operator
type NotContainsInsensitiveOp OpType

// Name is not contains insensitive operator name
func (o NotContainsInsensitiveOp) Name() OpType {
	return NotContainsInsensitive
}

// ValidateValue validate not contains insensitive operator's value
func (o NotContainsInsensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not contains insensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not contains insensitive operator's field and value to a mongo query condition.
func (o NotContainsInsensitiveOp) ToMgo(
	field string,
	value interface{},
) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBNot: map[string]interface{}{
				DBLIKE:    value,
				DBOPTIONS: "i",
			},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotContainsInsensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.Contains(strings.ToLower(val1), strings.ToLower(val2)), nil
}
