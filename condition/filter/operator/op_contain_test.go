package operator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsValidate(t *testing.T) {
	op := GetOperator(Contains)

	// test contains string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid contains type
	err = op.ValidateValue("", nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(1, nil)
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

func TestContainsMongoCond(t *testing.T) {
	op := GetOperator(Contains)

	// test contains string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBLIKE:    "a",
		DBOPTIONS: "i",
	}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestContainsSensitiveValidate(t *testing.T) {
	op := GetOperator(ContainsSensitive)

	// test contains string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid contains type
	err = op.ValidateValue("", nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(1, nil)
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

func TestContainsSensitiveMongoCond(t *testing.T) {
	op := GetOperator(ContainsSensitive)

	// test contains string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBLIKE: "a"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestNotContainsValidate(t *testing.T) {
	op := GetOperator(NotContains)

	// test not contains string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid not contains type
	err = op.ValidateValue("", nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(1, nil)
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

func TestNotContainsMongoCond(t *testing.T) {
	op := GetOperator(NotContains)

	// test not contains string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBNot: map[string]interface{}{DBLIKE: "a"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestNotContainsInsensitiveValidate(t *testing.T) {
	op := GetOperator(NotContainsInsensitive)

	// test not contains insensitive string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid not contains insensitive type
	err = op.ValidateValue("", nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(1, nil)
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

func TestNotContainsInsensitiveMongoCond(t *testing.T) {
	op := GetOperator(NotContainsInsensitive)

	// test not contains insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBNot: map[string]interface{}{
			DBLIKE: "a", DBOPTIONS: "i",
		}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestContainsMatch(t *testing.T) {
	op := GetOperator(Contains)

	// test contains string type
	matched, err := op.Match("123aBcdef", "Ab")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestContainsSensitiveMatch(t *testing.T) {
	op := GetOperator(ContainsSensitive)

	// test contains string type
	matched, err := op.Match("123abcdef", "ab")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match("123abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotContainsMatch(t *testing.T) {
	op := GetOperator(NotContains)

	// test not contains string type
	matched, err := op.Match("123abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "ab")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotContainsInsensitiveMatch(t *testing.T) {
	op := GetOperator(NotContainsInsensitive)

	// test not contains insensitive string type
	matched, err := op.Match("123abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123Abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
