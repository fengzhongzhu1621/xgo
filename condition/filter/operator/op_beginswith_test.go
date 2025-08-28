package operator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeginsWithValidate(t *testing.T) {
	op := GetOperator(BeginsWith)

	// test begins with string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid begins with type
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

func TestBeginsWithMongoCond(t *testing.T) {
	op := GetOperator(BeginsWith)

	// test begins with string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBLIKE: "^a",
	}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestBeginsWithInsensitiveValidate(t *testing.T) {
	op := GetOperator(BeginsWithInsensitive)

	// test begins with insensitive string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid begins with insensitive type
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

func TestBeginsWithInsensitiveMongoCond(t *testing.T) {
	op := GetOperator(BeginsWithInsensitive)

	// test begins with insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBLIKE:    "^a",
		DBOPTIONS: "i",
	}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestNotBeginsWithValidate(t *testing.T) {
	op := GetOperator(NotBeginsWith)

	// test not begins with string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid not begins with type
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

func TestNotBeginsWithMongoCond(t *testing.T) {
	op := GetOperator(NotBeginsWith)

	// test not begins with string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBNot: map[string]interface{}{DBLIKE: "^a"},
	}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestNotBeginsWithInsensitiveValidate(t *testing.T) {
	op := GetOperator(NotBeginsWithInsensitive)

	// test not begins with insensitive string type
	err := op.ValidateValue("a", nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid not begins with insensitive type
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

func TestNotBeginsWithInsensitiveMongoCond(t *testing.T) {
	op := GetOperator(NotBeginsWithInsensitive)

	// test not begins with insensitive string type
	cond, err := op.ToMgo("test", "a")
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBNot: map[string]interface{}{DBLIKE: "^a", DBOPTIONS: "i"},
	}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestBeginsWithMatch(t *testing.T) {
	op := GetOperator(BeginsWith)

	// test begins with string type
	matched, err := op.Match("abcdef", "ab")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match("abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestBeginsWithInsensitiveMatch(t *testing.T) {
	op := GetOperator(BeginsWithInsensitive)

	// test begins with insensitive string type
	matched, err := op.Match("aBcdef", "Ab")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("Abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotBeginsWithMatch(t *testing.T) {
	op := GetOperator(NotBeginsWith)

	// test not begins with string type
	matched, err := op.Match("abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("abcdef", "aB")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("abcdef", "ab")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func TestNotBeginsWithInsensitiveMatch(t *testing.T) {
	op := GetOperator(NotBeginsWithInsensitive)

	// test not begins with insensitive string type
	matched, err := op.Match("abcdef", "ac")
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match("aBcdef", "Ab")
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
