package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// ValidateSliceOfBasicType validate if the value is a slice of basic type
func ValidateSliceOfBasicType(value interface{}, limit uint) error {
	if value == nil {
		return errors.New("value is nil")
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return fmt.Errorf("value(%+v) is not of array type", value)
	}

	v := reflect.ValueOf(value)
	length := v.Len()
	if length == 0 {
		return errors.New("value is empty")
	}

	if length > int(limit) {
		return fmt.Errorf("elements length %d exceeds maximum %d", length, limit)
	}

	// each element in the array or slice should be of the same basic type.
	var firstItemType string
	for i := 0; i < length; i++ {
		item := v.Index(i).Interface()

		var itemType string
		switch item.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, json.Number:
			itemType = "numeric"
		case bool:
			itemType = "bool"
		case string:
			itemType = "string"
		default:
			return fmt.Errorf("array element index(%d) value(%+v) is not of basic type", i, item)
		}

		// 判断第一个元素的类型
		if i == 0 {
			firstItemType = itemType
			continue
		}

		if firstItemType != itemType {
			return fmt.Errorf("array element index(%d) value(%+v) type is not %s", i, item, firstItemType)
		}
	}

	return nil
}
