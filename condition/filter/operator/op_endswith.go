package operator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// EndsWithOp is ends with operator
type EndsWithOp OpType

// Name is ends with operator name
func (o EndsWithOp) Name() OpType {
	return EndsWith
}

// ValidateValue validate ends with operator's value
func (o EndsWithOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("ends with operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the ends with operator's field and value to a mongo query condition.
func (o EndsWithOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBLIKE: fmt.Sprintf("%s$", value),
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o EndsWithOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.HasSuffix(val1, val2), nil
}

// EndsWithInsensitiveOp is ends with insensitive operator
type EndsWithInsensitiveOp OpType

// Name is ends with insensitive operator name
func (o EndsWithInsensitiveOp) Name() OpType {
	return EndsWithInsensitive
}

// ValidateValue validate ends with insensitive operator's value
func (o EndsWithInsensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("ends with insensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the ends with insensitive operator's field and value to a mongo query condition.
func (o EndsWithInsensitiveOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBLIKE:    fmt.Sprintf("%s$", value),
			DBOPTIONS: "i",
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o EndsWithInsensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.HasSuffix(strings.ToLower(val1), strings.ToLower(val2)), nil
}

// NotEndsWithOp is not ends with operator
type NotEndsWithOp OpType

// Name is not ends with operator name
func (o NotEndsWithOp) Name() OpType {
	return NotEndsWith
}

// ValidateValue validate not ends with operator's value
func (o NotEndsWithOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not ends with operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not ends with operator's field and value to a mongo query condition.
func (o NotEndsWithOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBNot: map[string]interface{}{DBLIKE: fmt.Sprintf("%s$", value)},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotEndsWithOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.HasSuffix(val1, val2), nil
}

// NotEndsWithInsensitiveOp is not ends with insensitive operator
type NotEndsWithInsensitiveOp OpType

// Name is not ends with insensitive operator name
func (o NotEndsWithInsensitiveOp) Name() OpType {
	return NotEndsWithInsensitive
}

// ValidateValue validate not ends with insensitive operator's value
func (o NotEndsWithInsensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not ends with insensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// Match checks if the first data matches the second data by this operator
func (o NotEndsWithInsensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.HasSuffix(strings.ToLower(val1), strings.ToLower(val2)), nil
}

// ToMgo convert the not ends with insensitive operator's field and value to a mongo query condition.
func (o NotEndsWithInsensitiveOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBNot: map[string]interface{}{
				DBLIKE:    fmt.Sprintf("%s$", value),
				DBOPTIONS: "i",
			},
		},
	}, nil
}
