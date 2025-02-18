package operator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndsWithValidate(t *testing.T) {
	op := GetOperator(EndsWith)

	// test ends with string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid ends with type
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

func TestEndsWithMongoCond(t *testing.T) {
	op := GetOperator(EndsWith)

	// test ends with string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBLIKE: "a$"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestEndsWithInsensitiveValidate(t *testing.T) {
	op := GetOperator(EndsWithInsensitive)

	// test ends with insensitive string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid ends with insensitive type
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

func TestEndsWithInsensitiveMongoCond(t *testing.T) {
	op := GetOperator(EndsWithInsensitive)

	// test ends with insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBLIKE: "a$",
		DBOPTIONS: "i"}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestNotEndsWithValidate(t *testing.T) {
	op := GetOperator(NotEndsWith)

	// test not ends with string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid not ends with type
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

func TestNotEndsWithMongoCond(t *testing.T) {
	op := GetOperator(NotEndsWith)

	// test not ends with string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBNot: map[string]interface{}{DBLIKE: "a$"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestNotEndsWithInsensitiveValidate(t *testing.T) {
	op := GetOperator(NotEndsWithInsensitive)

	// test not ends with insensitive string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid not ends with insensitive type
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

func TestNotEndsWithInsensitiveMongoCond(t *testing.T) {
	op := GetOperator(NotEndsWithInsensitive)

	// test not ends with insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBNot: map[string]interface{}{DBLIKE: "a$", DBOPTIONS: "i"}}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestEndsWithMatch(t *testing.T) {
	op := GetOperator(EndsWith)

	// test ends with string type
	matched, err := op.Match("123abcdef", "ef")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "eF")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match("123abcdef", "df")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestEndsWithInsensitiveMatch(t *testing.T) {
	op := GetOperator(EndsWithInsensitive)

	// test ends with insensitive string type
	matched, err := op.Match("123abcDef", "dEf")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123Abcdef", "abc")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotEndsWithMatch(t *testing.T) {
	op := GetOperator(NotEndsWith)

	// test not ends with string type
	matched, err := op.Match("123abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcdef", "ef")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotEndsWithInsensitiveMatch(t *testing.T) {
	op := GetOperator(NotEndsWithInsensitive)

	// test not ends with insensitive string type
	matched, err := op.Match("123Abcdef", "abc")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("123abcDef", "dEf")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
