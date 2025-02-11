package operator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// EqualOp is equal operator type
type EqualOp OpType

// Name is equal operator name
func (o EqualOp) Name() OpType {
	return Equal
}

// ValidateValue validate equal operator's value
func (o EqualOp) ValidateValue(v interface{}, opt *ExprOption) error {
	// 判断匹配的值是否是内置基础类型
	if !validator.IsBasicValue(v) {
		return fmt.Errorf("invalid eq value(%+v)", v)
	}
	return nil
}

// ToMgo convert the equal operator's field and value to a mongo query condition.
func (o EqualOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	return mapstr.MapStr{
		field: map[string]interface{}{DBEQ: value},
	}, nil
}

// Match checks if the first data matches the second data by this operator
func (o EqualOp) Match(value1, value2 interface{}) (bool, error) {
	switch t := value1.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, json.Number:
		// interface{} 转换为 float64 类型进行比较
		val1, val2, err := parseNumericValues(value1, value2)
		if err != nil {
			return false, err
		}
		return val1 == val2, nil
	case string:
		val2, ok := value2.(string)
		if !ok {
			return false, fmt.Errorf("rule value type(%T) not matches input type(%T)", value2, value1)
		}
		return val2 == t, nil
	case bool:
		val2, ok := value2.(bool)
		if !ok {
			return false, fmt.Errorf("rule value type(%T) not matches input type(%T)", value2, value1)
		}
		return val2 == t, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("value(%+v) is not of basic type", value1)
	}
}
