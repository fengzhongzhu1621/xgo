package operator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualValidate(t *testing.T) {
	op := GetOperator(Equal)

	// test equal int type
	err := op.ValidateValue(1, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test equal string type
	err = op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test equal bool type
	err = op.ValidateValue(false, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid equal type
	err = op.ValidateValue([]int64{1}, nil)
	if err == nil {
		t.Errorf("to mongo should return error")
		return
	}

	err = op.ValidateValue(map[string]interface{}{"test1": 1}, nil)
	if err == nil {
		t.Errorf("to mongo should return error")
		return
	}

	err = op.ValidateValue(struct{}{}, nil)
	if err == nil {
		t.Errorf("to mongo should return error")
		return
	}

	err = op.ValidateValue(nil, nil)
	if err == nil {
		t.Errorf("to mongo should return error")
		return
	}
}

func TestEqualMongoCond(t *testing.T) {
	op := GetOperator(Equal)

	// test equal int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBEQ: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test equal string type
	cond, err = op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBEQ: "a"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test equal bool type
	cond, err = op.ToMgo("test", false)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(
		cond,
		map[string]interface{}{"test": map[string]interface{}{DBEQ: false}},
	) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestEqualMatch(t *testing.T) {
	op := GetOperator(Equal)

	// test equal int type
	matched, err := op.Match(int32(1), 1.0)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(int32(2), 1.0)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test equal string type
	matched, err = op.Match("a", "a")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("a", "b")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test equal bool type
	matched, err = op.Match(false, false)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(true, false)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
