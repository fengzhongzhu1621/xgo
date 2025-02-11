package operator

import (
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/cast"
)

// parseNumericValues interface{} 转换为 float64 类型
func parseNumericValues(value1, value2 interface{}) (float64, float64, error) {
	val1, err := cast.ToFloat64E(value1)
	if err != nil {
		return 0, 0, fmt.Errorf("parse input value(%+v) failed, err: %v", value1, err)
	}

	val2, err := cast.ToFloat64E(value2)
	if err != nil {
		return 0, 0, fmt.Errorf("parse rule value(%+v) failed, err: %v", value2, err)
	}

	return val1, val2, nil
}

func parseStringValues(value1, value2 interface{}) (string, string, error) {
	val1, ok := value1.(string)
	if !ok {
		return "", "", fmt.Errorf("input value(%+v) is not string type", value1)
	}

	val2, ok := value2.(string)
	if !ok {
		return "", "", fmt.Errorf("rule value(%+v) is not string type", value2)
	}

	return val1, val2, nil
}

func parseTimeValues(value1, value2 interface{}) (time.Time, time.Time, error) {
	val1, err := cast.ConvToTime(value1)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse input value(%+v) failed, err: %v", value1, err)
	}

	val2, err := cast.ConvToTime(value2)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse rule value(%+v) failed, err: %v", value2, err)
	}

	return val1, val2, nil
}
