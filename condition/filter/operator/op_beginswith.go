package operator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// BeginsWithOp is begins with operator
type BeginsWithOp OpType

// Name is begins with operator name
func (o BeginsWithOp) Name() OpType {
	return BeginsWith
}

// ValidateValue validate begins with operator's value
func (o BeginsWithOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("begins with operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the begins with operator's field and value to a mongo query condition.
func (o BeginsWithOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBLIKE: fmt.Sprintf("^%s", value),
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o BeginsWithOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(val1, val2), nil
}

// BeginsWithInsensitiveOp is begins with insensitive operator
type BeginsWithInsensitiveOp OpType

// Name is begins with insensitive operator name
func (o BeginsWithInsensitiveOp) Name() OpType {
	return BeginsWithInsensitive
}

// ValidateValue validate begins with insensitive operator's value
func (o BeginsWithInsensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("begins with insensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the begins with insensitive operator's field and value to a mongo query condition.
func (o BeginsWithInsensitiveOp) ToMgo(
	field string,
	value interface{},
) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBLIKE:    fmt.Sprintf("^%s", value),
			DBOPTIONS: "i",
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o BeginsWithInsensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(strings.ToLower(val1), strings.ToLower(val2)), nil
}

// NotBeginsWithOp is not begins with operator
type NotBeginsWithOp OpType

// Name is not begins with operator name
func (o NotBeginsWithOp) Name() OpType {
	return NotBeginsWith
}

// ValidateValue validate not begins with operator's value
func (o NotBeginsWithOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not begins with operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not begins with operator's field and value to a mongo query condition.
func (o NotBeginsWithOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBNot: map[string]interface{}{
				DBLIKE: fmt.Sprintf("^%s", value),
			},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotBeginsWithOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.HasPrefix(val1, val2), nil
}

// NotBeginsWithInsensitiveOp is not begins with insensitive operator
type NotBeginsWithInsensitiveOp OpType

// Name is not begins with insensitive operator name
func (o NotBeginsWithInsensitiveOp) Name() OpType {
	return NotBeginsWithInsensitive
}

// ValidateValue validate not begins with insensitive operator's value
func (o NotBeginsWithInsensitiveOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateNotEmptyStringType(v)
	if err != nil {
		return fmt.Errorf("not begins with insensitive operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not begins with insensitive operator's field and value to a mongo query condition.
func (o NotBeginsWithInsensitiveOp) ToMgo(
	field string,
	value interface{},
) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{
			DBNot: map[string]interface{}{
				DBLIKE:    fmt.Sprintf("^%s", value),
				DBOPTIONS: "i",
			},
		},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotBeginsWithInsensitiveOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseStringValues(value1, value2)
	if err != nil {
		return false, err
	}
	return !strings.HasPrefix(strings.ToLower(val1), strings.ToLower(val2)), nil
}
