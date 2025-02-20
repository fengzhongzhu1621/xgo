package criteria

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/fengzhongzhu1621/xgo/validator"
)

// ValidateFieldValue 验证字段类型和需要匹配的值，即规则中的 Field 和 Value
func ValidateFieldValue(v interface{}, typ FieldType) error {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice:
		return ValidateSliceElements(v, typ)
	default:
	}

	switch typ {
	case String, Enum:
		if reflect.ValueOf(v).Type().Kind() != reflect.String {
			return errors.New("value should be a string")
		}

	case Numeric, Timestamp:
		if !validator.IsNumeric(v) {
			return errors.New("value should be a numeric")
		}

	case Boolean:
		if reflect.ValueOf(v).Type().Kind() != reflect.Bool {
			return errors.New("value should be a boolean")
		}

	case Time:
		if err := validator.ValidateDatetimeType(v); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported value type format: %s", typ)
	}

	return nil
}

func ValidateSliceElements(v interface{}, typ FieldType) error {
	value := reflect.ValueOf(v)
	length := value.Len()
	if length == 0 {
		return nil
	}

	// validate each slice's element data type
	for i := 0; i < length; i++ {
		if err := ValidateFieldValue(value.Index(i).Interface(), typ); err != nil {
			return err
		}
	}

	return nil
}
