package operator

import (
	"errors"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// DatetimeLessOp is datetime less than operator
type DatetimeLessOp OpType

// Name is datetime less than operator name
func (o DatetimeLessOp) Name() OpType {
	return DatetimeLess
}

// ValidateValue validate datetime less than operator value
func (o DatetimeLessOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateDatetimeType(v)
	if err != nil {
		return fmt.Errorf("datetime less than operator's value is invalid, err: %v", err)
	}
	return nil
}

// ToMgo convert the datetime less than operator's field and value to a mongo query condition.
func (o DatetimeLessOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	v, err := cast.ConvToTime(value)
	if err != nil {
		return nil, fmt.Errorf("convert value to time failed, err: %v", err)
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBLT: v},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o DatetimeLessOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseTimeValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1.Before(val2), nil
}

// DatetimeLessOrEqualOp is datetime less than or equal operator
type DatetimeLessOrEqualOp OpType

// Name is datetime less than or equal operator name
func (o DatetimeLessOrEqualOp) Name() OpType {
	return DatetimeLessOrEqual
}

// ValidateValue validate datetime less than or equal operator value
func (o DatetimeLessOrEqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateDatetimeType(v)
	if err != nil {
		return fmt.Errorf("datetime less than or equal operator's value is invalid, err: %v", err)
	}
	return nil
}

// ToMgo convert the datetime less than or equal operator's field and value to a mongo query condition.
func (o DatetimeLessOrEqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	v, err := cast.ConvToTime(value)
	if err != nil {
		return nil, fmt.Errorf("convert value to time failed, err: %v", err)
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBLTE: v},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o DatetimeLessOrEqualOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseTimeValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1.Before(val2) || val1.Equal(val2), nil
}

// DatetimeGreaterOp is datetime greater than operator
type DatetimeGreaterOp OpType

// Name is datetime greater than operator name
func (o DatetimeGreaterOp) Name() OpType {
	return DatetimeGreater
}

// ValidateValue validate datetime greater than operator value
func (o DatetimeGreaterOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateDatetimeType(v)
	if err != nil {
		return fmt.Errorf("datetime greater than operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the datetime greater than operator's field and value to a mongo query condition.
func (o DatetimeGreaterOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	v, err := cast.ConvToTime(value)
	if err != nil {
		return nil, fmt.Errorf("convert value to time failed, err: %v", err)
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBGT: v},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o DatetimeGreaterOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseTimeValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1.After(val2), nil
}

// DatetimeGreaterOrEqualOp is datetime greater than or equal operator
type DatetimeGreaterOrEqualOp OpType

// Name is datetime greater than or equal operator name
func (o DatetimeGreaterOrEqualOp) Name() OpType {
	return DatetimeGreaterOrEqual
}

// ValidateValue validate datetime greater than or equal operator value
func (o DatetimeGreaterOrEqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	err := validator.ValidateDatetimeType(v)
	if err != nil {
		return fmt.Errorf("datetime greater than or equal operator's value is invalid, err: %v", err)
	}
	return nil
}

// ToMgo convert the datetime greater than or equal operator's field and value to a mongo query condition.
func (o DatetimeGreaterOrEqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	v, err := cast.ConvToTime(value)
	if err != nil {
		return nil, fmt.Errorf("convert value to time failed, err: %v", err)
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBGTE: v},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o DatetimeGreaterOrEqualOp) Match(value1, value2 interface{}) (bool, error) {
	val1, val2, err := parseTimeValues(value1, value2)
	if err != nil {
		return false, err
	}
	return val1.After(val2) || val1.Equal(val2), nil
}
