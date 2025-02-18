package operator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLessOrEqualValidate(t *testing.T) {
	op := GetOperator(LessOrEqual)

	// test less or equal int type
	err := op.ValidateValue(1, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test less or equal than 0
	err = op.ValidateValue(uint64(0), nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test less or equal than a negative number
	err = op.ValidateValue(int32(-1), nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid less or equal type
	err = op.ValidateValue("a", nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(false, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue([]int64{1}, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(map[string]interface{}{"test1": 1}, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(struct{}{}, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(nil, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}
}

func TestLessOrEqualMongoCond(t *testing.T) {
	op := GetOperator(LessOrEqual)

	// test less or equal int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBLTE: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test less or equal than 0
	cond, err = op.ToMgo("test", uint64(0))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBLTE: uint64(0)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test less or equal than a negative number
	cond, err = op.ToMgo("test", int32(-1))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBLTE: int32(-1)}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestLessOrEqualMatch(t *testing.T) {
	op := GetOperator(LessOrEqual)

	// test less or equal int type
	matched, err := op.Match(0.01, 1)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(3, 1)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test less or equal than 0
	matched, err = op.Match(-1, uint64(0))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(1.1, uint64(0))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test less or equal than a negative number
	matched, err = op.Match(-1.0, int32(-1))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(-0.01, int32(-1))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
