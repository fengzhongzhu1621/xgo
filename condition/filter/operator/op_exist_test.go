package operator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExistValidate(t *testing.T) {
	op := GetOperator(Exist)

	// test exist cond
	err := op.ValidateValue(nil, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}
}

func TestExistMongoCond(t *testing.T) {
	op := GetOperator(Exist)

	// test exist cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(
		cond,
		map[string]interface{}{"test": map[string]interface{}{DBExists: true}},
	) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestNotExistValidate(t *testing.T) {
	op := GetOperator(NotExist)

	// test not exist cond
	err := op.ValidateValue(nil, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}
}

func TestNotExistMongoCond(t *testing.T) {
	op := GetOperator(NotExist)

	// test not exist cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(
		cond,
		map[string]interface{}{"test": map[string]interface{}{DBExists: false}},
	) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestExistMatch(t *testing.T) {
	op := GetOperator(Exist)

	// test exist matched
	matched, err := op.Match(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(1, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotExistMatch(t *testing.T) {
	op := GetOperator(NotExist)

	// test not exist matched
	matched, err := op.Match(1, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
