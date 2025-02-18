package operator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
	"github.com/stretchr/testify/assert"
)

// NotInOp is not in operator
type NotInOp OpType

// Name is not in operator name
func (o NotInOp) Name() OpType {
	return NotIn
}

// ValidateValue validate not in value
func (o NotInOp) ValidateValue(v interface{}, opt *ExprOption) error {
	if opt == nil {
		return errors.New("validate option must be set")
	}

	err := validator.ValidateSliceOfBasicType(v, opt.MaxNotInLimit)
	if err != nil {
		return fmt.Errorf("nin operator's value is invalid, err: %v", err)
	}

	return nil
}

// ToMgo convert the not in operator's field and value to a mongo query condition.
func (o NotInOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBNIN: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o NotInOp) Match(value1, value2 interface{}) (bool, error) {
	matched, err := GetOperator(In).Match(value1, value2)
	if err != nil {
		return false, err
	}
	return !matched, nil
}

func TestNotInMatch(t *testing.T) {
	op := GetOperator(NotIn)

	// test not in int array type
	matched, err := op.Match(1.0, []int64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match(3, []int64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	// test not in string array type
	matched, err = op.Match("b", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match("c", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	// test not in bool array type
	matched, err = op.Match(true, []interface{}{false, true})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match(false, []interface{}{true})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)
}
