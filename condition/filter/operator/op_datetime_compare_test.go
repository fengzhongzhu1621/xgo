package operator

import (
	"reflect"
	"testing"
	"time"
)

func TestDatetimeLessValidate(t *testing.T) {
	op := GetOperator(DatetimeLess)

	// test datetime less time type
	now := time.Now()
	err := op.ValidateValue(now, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test datetime less timestamp type
	err = op.ValidateValue(now.Unix(), nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test datetime less time string
	nowStr := now.Format(time.DateTime)
	err = op.ValidateValue(nowStr, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid datetime less type
	err = op.ValidateValue("2022", nil)
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

func TestDatetimeLessMongoCond(t *testing.T) {
	op := GetOperator(DatetimeLess)

	// test datetime less time type
	now := time.Now()
	cond, err := op.ToMgo("test", now)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBLT: now}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime less timestamp type
	nowUnixTime := time.Unix(now.Unix(), 0)
	cond, err = op.ToMgo("test", now.Unix())
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBLT: nowUnixTime}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime less time string
	nowStr := now.Format(time.DateTime)
	nowStrTime, _ := time.ParseInLocation(time.DateTime, nowStr, time.Local)
	cond, err = op.ToMgo("test", nowStr)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBLT: nowStrTime.UTC()}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestDatetimeLessOrEqualValidate(t *testing.T) {
	op := GetOperator(DatetimeLessOrEqual)

	// test datetime less or equal time type
	now := time.Now()
	err := op.ValidateValue(now, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test datetime less or equal timestamp type
	err = op.ValidateValue(now.Unix(), nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test datetime less or equal time string
	nowStr := now.Format(time.DateTime)
	err = op.ValidateValue(nowStr, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid datetime less or equal type
	err = op.ValidateValue("2022", nil)
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

func TestDatetimeLessOrEqualMongoCond(t *testing.T) {
	op := GetOperator(DatetimeLessOrEqual)

	// test datetime less or equal time type
	now := time.Now()
	cond, err := op.ToMgo("test", now)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBLTE: now}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime less or equal timestamp type
	nowUnixTime := time.Unix(now.Unix(), 0)
	cond, err = op.ToMgo("test", now.Unix())
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBLTE: nowUnixTime}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime less or equal time string
	nowStr := now.Format(time.DateTime)
	nowStrTime, _ := time.ParseInLocation(time.DateTime, nowStr, time.Local)
	cond, err = op.ToMgo("test", nowStr)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBLTE: nowStrTime.UTC()}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestDatetimeGreaterValidate(t *testing.T) {
	op := GetOperator(DatetimeGreater)

	// test datetime greater time type
	now := time.Now()
	err := op.ValidateValue(now, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test datetime greater timestamp type
	err = op.ValidateValue(now.Unix(), nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test datetime greater time string
	nowStr := now.Format(time.DateTime)
	err = op.ValidateValue(nowStr, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid datetime greater type
	err = op.ValidateValue("2022", nil)
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

func TestDatetimeGreaterMongoCond(t *testing.T) {
	op := GetOperator(DatetimeGreater)

	// test datetime greater time type
	now := time.Now()
	cond, err := op.ToMgo("test", now)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBGT: now}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime greater timestamp type
	nowUnixTime := time.Unix(now.Unix(), 0)
	cond, err = op.ToMgo("test", now.Unix())
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBGT: nowUnixTime}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime greater time string
	nowStr := now.Format(time.DateTime)
	nowStrTime, _ := time.ParseInLocation(time.DateTime, nowStr, time.Local)
	cond, err = op.ToMgo("test", nowStr)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBGT: nowStrTime.UTC()}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestDatetimeGreaterOrEqualValidate(t *testing.T) {
	op := GetOperator(DatetimeGreaterOrEqual)

	// test datetime greater or equal time type
	now := time.Now()
	err := op.ValidateValue(now, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test datetime greater or equal timestamp type
	err = op.ValidateValue(now.Unix(), nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test datetime greater or equal time string
	nowStr := now.Format(time.DateTime)
	err = op.ValidateValue(nowStr, nil)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid datetime greater or equal type
	err = op.ValidateValue("2022", nil)
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

func TestDatetimeGreaterOrEqualMongoCond(t *testing.T) {
	op := GetOperator(DatetimeGreaterOrEqual)

	// test datetime greater or equal time type
	now := time.Now()
	cond, err := op.ToMgo("test", now)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBGTE: now}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime greater or equal timestamp type
	nowUnixTime := time.Unix(now.Unix(), 0)
	cond, err = op.ToMgo("test", now.Unix())
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{DBGTE: nowUnixTime}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}

	// test datetime greater or equal time string
	nowStr := now.Format(time.DateTime)
	nowStrTime, _ := time.ParseInLocation(time.DateTime, nowStr, time.Local)
	cond, err = op.ToMgo("test", nowStr)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(cond, map[string]interface{}{"test": map[string]interface{}{
		DBGTE: nowStrTime.UTC()}}) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}
