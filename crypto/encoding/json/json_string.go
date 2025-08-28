package json

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// JsonString is the json string type
type JsonString string

// GetValue get value by field from json string
func (j JsonString) GetValue(field string) (interface{}, error) {
	if j == "" {
		return nil, fmt.Errorf("json data is empty")
	}
	val := gjson.Get(string(j), field)
	switch val.Type {
	case gjson.Null:
		return nil, nil
	case gjson.Number:
		if strings.Contains(val.Raw, ".") {
			return val.Float(), nil
		}
		return val.Int(), nil
	default:
		return val.Value(), nil
	}
}
