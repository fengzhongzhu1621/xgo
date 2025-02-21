package operator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInValidate(t *testing.T) {
	op := GetOperator(In)

	// test in int array type
	err := op.ValidateValue([]int64{1, 2}, NewDefaultExprOpt(nil))
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test in string array type
	err = op.ValidateValue([]string{"a", "b"}, NewDefaultExprOpt(nil))
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test in bool array type
	err = op.ValidateValue([]interface{}{true, false}, NewDefaultExprOpt(nil))
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid in type
	err = op.ValidateValue(1, NewDefaultExprOpt(nil))
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue("a", NewDefaultExprOpt(nil))
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(false, NewDefaultExprOpt(nil))
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(map[string]interface{}{"test1": 1}, NewDefaultExprOpt(nil))
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(struct{}{}, NewDefaultExprOpt(nil))
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(nil, NewDefaultExprOpt(nil))
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue([]int64{}, NewDefaultExprOpt(nil))
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue([]interface{}{1, "a"}, NewDefaultExprOpt(nil))
	if err == nil {
		t.Errorf("validate should return error")
		return
	}
}

func TestInMongoCond(t *testing.T) {
	op := GetOperator(In)

	// test in int array type
	cond, err := op.ToMgo("test", []int64{1, 2})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBIN: []int64{1, 2}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test in string array type
	cond, err = op.ToMgo("test", []string{"a", "b"})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBIN: []string{"a", "b"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test in bool array type
	cond, err = op.ToMgo("test", []interface{}{true, false})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBIN: []interface{}{true, false}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestInMatch(t *testing.T) {
	op := GetOperator(In)

	// test in int array type
	matched, err := op.Match(1.0, []int64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(3, []int64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test in string array type
	matched, err = op.Match("b", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("c", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	// test in bool array type
	matched, err = op.Match(false, []interface{}{true, false})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match(false, []interface{}{true})
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
