package operator

import (
	"reflect"
	"testing"
)

func TestNotInValidate(t *testing.T) {
	op := GetOperator(NotIn)

	// test not in int array type
	err := op.ValidateValue([]int64{1, 2}, NewDefaultExprOpt(nil))
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test not in string array type
	err = op.ValidateValue([]string{"a", "b"}, NewDefaultExprOpt(nil))
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test not in bool array type
	err = op.ValidateValue([]interface{}{true, false}, NewDefaultExprOpt(nil))
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid not in type
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

func TestNotInMongoCond(t *testing.T) {
	op := GetOperator(NotIn)

	// test not in int array type
	cond, err := op.ToMgo("test", []int64{1, 2})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBNIN: []int64{1, 2}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test not in string array type
	cond, err = op.ToMgo("test", []string{"a", "b"})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBNIN: []string{"a", "b"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test not in bool array type
	cond, err = op.ToMgo("test", []interface{}{true, false})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBNIN: []interface{}{true, false}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}
