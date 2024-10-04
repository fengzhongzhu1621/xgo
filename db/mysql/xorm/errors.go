package xorm

import "strings"

func IsDataExists(err error) bool {
	if err == nil {
		return true
	}

	if strings.Contains(err.Error(), "1062") && strings.Contains(err.Error(), "Duplicate") {
		return true
	}

	return false
}
