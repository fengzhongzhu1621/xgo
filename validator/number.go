package validator

import (
	"encoding/json"
	"regexp"
)

const (
	numCharPattern = `^[a-zA-Z0-9]*$`
)

var (
	numCharRegexp = regexp.MustCompile(numCharPattern)
)

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

// IsNumChar 是否字母、数字组合
func IsNumChar(sInput string) bool {
	return numCharRegexp.MatchString(sInput)
}
