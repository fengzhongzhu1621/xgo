package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

const (
	charPattern = `^[a-zA-Z]*$`
)

var charRegexp = regexp.MustCompile(charPattern)

// IsLower 判断字符串是否包含小写字母.
func IsLower(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

// ValidateStringType validate if the value is a string type
func ValidateStringType(value interface{}) error {
	if reflect.TypeOf(value).Kind() != reflect.String {
		return fmt.Errorf("value(%+v) is not of string type", value)
	}
	return nil
}

// ValidateNotEmptyStringType validate if the value is a not empty string type
func ValidateNotEmptyStringType(value interface{}) error {
	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("value(%+v) is not of string type", value)
	}

	if len(strVal) == 0 {
		return errors.New("value is empty")
	}
	return nil
}

// IsChar 是否大、小写字母组合
func IsChar(sInput string) bool {
	return charRegexp.MatchString(sInput)
}

// CheckLen 字符串输入长度
func CheckLen(sInput string, min, max int) bool {
	if len(sInput) >= min && len(sInput) <= max {
		return true
	}
	return false
}
