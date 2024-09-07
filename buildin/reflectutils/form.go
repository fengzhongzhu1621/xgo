package reflectutils

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	json "github.com/fengzhongzhu1621/xgo/crypto/encoding/json"
	bytesconv "github.com/fengzhongzhu1621/xgo/str/bytesconv"
)

var ErrUnknownType = errors.New("unknown type")

func SetTimeDuration(val string, value reflect.Value, field reflect.StructField) error {
	// 获得时间间隔
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

func SetIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func SetUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func SetBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func SetFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func SetTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	// 获得日期字段的格式
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	// 判断是否是时间戳
	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

func SetWithProperType(val string, value reflect.Value, field reflect.StructField) error {
	switch value.Kind() {
	case reflect.Int:
		return SetIntField(val, 0, value)
	case reflect.Int8:
		return SetIntField(val, 8, value)
	case reflect.Int16:
		return SetIntField(val, 16, value)
	case reflect.Int32:
		return SetIntField(val, 32, value)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return SetTimeDuration(val, value, field)
		}
		return SetIntField(val, 64, value)
	case reflect.Uint:
		return SetUintField(val, 0, value)
	case reflect.Uint8:
		return SetUintField(val, 8, value)
	case reflect.Uint16:
		return SetUintField(val, 16, value)
	case reflect.Uint32:
		return SetUintField(val, 32, value)
	case reflect.Uint64:
		return SetUintField(val, 64, value)
	case reflect.Bool:
		return SetBoolField(val, value)
	case reflect.Float32:
		return SetFloatField(val, 32, value)
	case reflect.Float64:
		return SetFloatField(val, 64, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return SetTimeField(val, field, value)
		}
		return json.Unmarshal(bytesconv.StringToBytes(val), value.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal(bytesconv.StringToBytes(val), value.Addr().Interface())
	default:
		return ErrUnknownType
	}
	return nil
}

func SetArray(vals []string, value reflect.Value, field reflect.StructField) error {
	for i, s := range vals {
		err := SetWithProperType(s, value.Index(i), field)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetSlice(vals []string, value reflect.Value, field reflect.StructField) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := SetArray(vals, slice, field)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}
