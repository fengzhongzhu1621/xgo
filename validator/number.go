package validator

import "encoding/json"

// IsNumeric judges if value is a number
func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, json.Number:
		return true
	}

	return false
}

// IsInteger judges if value is a integer
func IsInteger(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, json.Number:
		return true
	}

	return false
}
