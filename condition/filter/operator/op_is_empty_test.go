package operator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmptyValidate(t *testing.T) {
	op := GetOperator(IsEmpty)

	// test is empty cond
	err := op.ValidateValue(nil, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}
}

func TestIsEmptyMongoCond(t *testing.T) {
	op := GetOperator(IsEmpty)

	// test is empty cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBSize: 0}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestIsNotEmptyValidate(t *testing.T) {
	op := GetOperator(IsNotEmpty)

	// test is not empty cond
	err := op.ValidateValue(nil, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}
}

func TestIsNotEmptyMongoCond(t *testing.T) {
	op := GetOperator(IsNotEmpty)

	// test is not empty cond
	cond, err := op.ToMgo("test", nil)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBSize: map[string]interface{}{DBGT: 0}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestIsEmptyMatch(t *testing.T) {
	op := GetOperator(IsEmpty)

	// test is empty matched
	matched, err := op.Match([]int64{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]float64{1, 2}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]string{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]string{""}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]bool{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]bool{true}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestIsNotEmptyMatch(t *testing.T) {
	op := GetOperator(IsNotEmpty)

	// test is not empty matched
	matched, err := op.Match([]int64{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]float64{1, 2}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]string{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]string{""}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]bool{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]bool{false}, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)
}
