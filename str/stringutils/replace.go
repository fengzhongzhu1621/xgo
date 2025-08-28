package stringutils

import "strings"

// IStringReplacer applies a set of replacements to a string.
type IStringReplacer interface {
	// Replace returns a copy of s with all replacements performed.
	Replace(s string) string
}

func RemoveStringSpace(data string) string {
	return strings.ReplaceAll(data, " ", "")
}

func PurifyStrings(ori string) string {
	purifyStr := strings.ReplaceAll(ori, "\\r", ";")
	purifyStr = strings.ReplaceAll(purifyStr, "\\n", ";")
	purifyStr = strings.ReplaceAll(purifyStr, "\n", ";")
	purifyStr = strings.ReplaceAll(purifyStr, "\r", ";")
	purifyStr = strings.ReplaceAll(purifyStr, " ", "")
	purifyStr = strings.ReplaceAll(purifyStr, ";", ",")

	return purifyStr
}

// RemoteTabCharacters 清理包含制表符和水平制表符的字符串
func RemoteTabCharacters(data string) string {
	data = strings.ReplaceAll(data, "\x09", "")
	data = strings.ReplaceAll(data, "\t", "")

	return data
}

// EncodeDot 将字符串中的点（.）编码为 Unicode 转义序列 \u002E
func EncodeDot(input string) string {
	return strings.ReplaceAll(input, ".", "\\u002E")
}

// DecodeDot 将字符串中的 Unicode 转义序列 \u002E 解码为点（.）
func DecodeDot(input string) string {
	return strings.ReplaceAll(input, "\\u002E", ".")
}
