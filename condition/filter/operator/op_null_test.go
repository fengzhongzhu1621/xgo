package operator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNullValidate(t *testing.T) {
	op := GetOperator(IsNull)

	// test is null cond
	err := op.ValidateValue(nil, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}
}

func TestIsNullMongoCond(t *testing.T) {
	op := GetOperator(IsNull)

	// test is null cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBEQ: nil}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestIsNotNullValidate(t *testing.T) {
	op := GetOperator(IsNotNull)

	// test is not null cond
	err := op.ValidateValue(nil, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}
}

func TestIsNotNullMongoCond(t *testing.T) {
	op := GetOperator(IsNotNull)

	// test is not null cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBNE: nil}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestIsNullMatch(t *testing.T) {
	op := GetOperator(IsNull)

	// test is null matched
	matched, err := op.Match(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(1, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestIsNotNullMatch(t *testing.T) {
	op := GetOperator(IsNotNull)

	// test is not null matched
	matched, err := op.Match(1, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
