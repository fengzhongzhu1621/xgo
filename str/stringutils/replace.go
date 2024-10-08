package stringutils

import "strings"

// IStringReplacer applies a set of replacements to a string.
type IStringReplacer interface {
	// Replace returns a copy of s with all replacements performed.
	Replace(s string) string
}

func RemoveStringSpace(data string) string {
	return strings.Replace(data, " ", "", -1)
}

func PurifyStrings(ori string) string {
	purifyStr := strings.Replace(ori, "\\r", ";", -1)
	purifyStr = strings.Replace(purifyStr, "\\n", ";", -1)
	purifyStr = strings.Replace(purifyStr, "\n", ";", -1)
	purifyStr = strings.Replace(purifyStr, "\r", ";", -1)
	purifyStr = strings.Replace(purifyStr, " ", "", -1)
	purifyStr = strings.Replace(purifyStr, ";", ",", -1)

	return purifyStr
}

// RemoteTabCharacters 清理包含制表符和水平制表符的字符串
func RemoteTabCharacters(data string) string {
	data = strings.Replace(data, "\x09", "", -1)
	data = strings.Replace(data, "\t", "", -1)

	return data
}
