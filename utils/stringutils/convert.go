package stringutils

import (
	"strconv"
)

func IsPrint(s string) bool {
	for _, c := range s {
		if !strconv.IsPrint(c) {
			return false
		}
	}

	return true
}

func QuoteIfNeeded(s string) string {
	if !IsPrint(s) {
		return strconv.Quote(s)
	}

	return s
}

func QuoteIfNeededV(s []string) []string {
	ret := make([]string, len(s))

	for i, v := range s {
		ret[i] = QuoteIfNeeded(v)
	}

	return ret
}

func QuoteV(s []string) []string {
	ret := make([]string, len(s))

	for i, v := range s {
		ret[i] = strconv.Quote(v)
	}

	return ret
}

func UnquoteIfPossible(s string) (string, error) {
	if len(s) == 0 || s[0] != '"' {
		return s, nil
	}

	return strconv.Unquote(s)
}
