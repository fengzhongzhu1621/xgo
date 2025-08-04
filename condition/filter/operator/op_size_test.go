package operator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizeValidate(t *testing.T) {
	op := GetOperator(Size)

	// test size int type
	err := op.ValidateValue(1, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test size equal to 0
	err = op.ValidateValue(uint64(0), nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid size type
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

	err = op.ValidateValue(int32(-1), nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}
}

func TestSizeMongoCond(t *testing.T) {
	op := GetOperator(Size)

	// test size int type
	cond, err := op.ToMgo("test", 1)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBSize: 1}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test size equal to 0
	cond, err = op.ToMgo("test", uint64(0))
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(
		cond,
		map[string]interface{}{"test": map[string]interface{}{DBSize: uint64(0)}},
	) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestSizeMatch(t *testing.T) {
	op := GetOperator(Size)

	// test size matched
	matched, err := op.Match([]int64{1, 2, 3}, 3)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]float64{1, 2}, 1)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]string{"1", "2"}, 2)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]string{""}, 2)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)

	matched, err = op.Match([]bool{}, 0)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = op.Match([]bool{true}, 2)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
