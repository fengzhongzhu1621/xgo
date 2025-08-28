package validator

import (
	"errors"
	"regexp"
)

const (
	// 小写字母开头, 可以包含小写字母/数字/下划线/连字符
	validIDString = "^[a-z]+[a-z0-9_-]*$"
)

var (
	ValidIDRegex = regexp.MustCompile(validIDString)

	ErrInvalidID = errors.New("invalid id: id should begin with a lowercase letter, " +
		"contains lowercase letters(a-z), numbers(0-9), underline(_) or hyphen(-)")
)
