package operator

import (
	"reflect"
	"testing"
)

func TestNotEqualValidate(t *testing.T) {
	op := GetOperator(NotEqual)

	// test not equal int type
	err := op.ValidateValue(1, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test not equal string type
	err = op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test not equal bool type
	err = op.ValidateValue(false, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid not equal type
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

func TestNotEqualMongoCond(t *testing.T) {
	op := GetOperator(NotEqual)

	// test not equal int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBNE: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test not equal string type
	cond, err = op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBNE: "a"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test not equal bool type
	cond, err = op.ToMgo("test", false)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}
	if !reflect.DeepEqual(
		cond,
		map[string]interface{}{"test": map[string]interface{}{DBNE: false}},
	) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}
