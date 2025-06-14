package stringutils

import (
	"fmt"
	"strings"
)

// Head 根据分隔符分割字符串，只分隔一次
func Head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}

func SplitString(r rune) bool {
	return r == ';' || r == ',' || r == '\n'
}

// SplitMulti split multi string by sep string.
// 不排除分隔后的空字符串
func SplitMulti(ss []string, sep string) []string {
	ns := make([]string, 0, len(ss)+1)
	for _, s := range ss {
		ns = append(ns, strings.Split(s, sep)...)
	}
	return ns
}

// ParseCommandlineMap parses a string in the format `key1:value1,key2:value2,...`
// where keys or values may contain commas and colons. We will allow escaping those using double quotes,
// so when passing in `"key1":"value1"`, we will not look inside the quoted sections.
// 解析类似 key1:value1,key2:value2,... 的输入字符串，并将其转换为 map[string]string
func ParseCommandlineMap(src string) (map[string]string, error) {
	result := make(map[string]string)

	// 按逗号分隔字符串，忽略逗号内的引号
	tuples := SplitStringRespectingQuotes(src, ',')
	for _, t := range tuples {
		// 分隔单个 key:value 对
		kv, err := splitKeyValue(t)
		if err != nil {
			return nil, err
		}
		result[kv.key] = kv.value
	}

	return result, nil
}

// splitKeyValue splits a single key-value pair into key and value, respecting quoted sections.
func splitKeyValue(pair string) (struct{ key, value string }, error) {
	var key strings.Builder
	var value strings.Builder
	inKeyQuotes := false   // 标记是否在 key 的引号内
	inValueQuotes := false // 标记是否在 value 的引号内
	escapeNext := false    // 标记下一个字符是否需要转义（如 \")
	inValue := false       // 标记是否已进入 value 部分（即是否已遇到第一个冒号）

	// 逐字符遍历当前 key:value 对
	for i, c := range pair {
		if escapeNext { // 如果当前字符需要转义（如前一个字符是 \）
			if inKeyQuotes { // 如果在 key 的引号内
				key.WriteRune(c)
			} else if inValueQuotes { // 如果在 value 的引号内
				value.WriteRune(c)
			} else { // 如果不在引号内（非法转义）
				return struct{ key, value string }{}, fmt.Errorf("invalid escape sequence at position %d", i)
			}
			escapeNext = false // 重置转义标记
			continue           // 跳过当前字符处理
		}

		switch c {
		case '\\': // 遇到转义符 \
			escapeNext = true // 标记下一个字符需要转义
		case '"':
			if inKeyQuotes {
				inKeyQuotes = false // 如果在 key 的引号内
			} else if inValueQuotes { // 如果在 value 的引号内
				inValueQuotes = false
			} else if inValue { // 如果已进入 value 部分
				// Value already started, so this must be closing quote
				inValueQuotes = true
			} else { // 如果尚未进入 value 部分
				// Key not started, so this must be opening quote
				inKeyQuotes = true
			}
		case ':':
			if inKeyQuotes { // 如果在 key 的引号内
				key.WriteRune(c) // 将冒号写入 key 缓冲区（视为普通字符）
			} else if inValueQuotes { // 如果在 value 的引号内
				value.WriteRune(c) // 将冒号写入 value 缓冲区（视为普通字符）
			} else if !inValue { // 如果尚未进入 value 部分
				// First colon outside quotes separates key and value
				inValue = true // 进入 value 部分（第一个冒号作为分隔符）
			} else { // 如果已在 value 部分
				// Colon inside value, just append
				value.WriteRune(c) // 将冒号写入 value 缓冲区（视为普通字符）
			}
		default: // 其他普通字符
			if !inValue { // 如果尚未进入 value 部分
				key.WriteRune(c) // 将字符写入 key 缓冲区
			} else { // 如果已进入 value 部分
				value.WriteRune(c) // 将字符写入 value 缓冲区
			}
		}
	}

	// 检查引号是否闭合
	// Validate that all quotes are closed
	if inKeyQuotes || inValueQuotes { // 如果 key 或 value 的引号未闭合
		return struct{ key, value string }{}, fmt.Errorf("unterminated quote in key-value pair")
	}

	// 返回解析结果（去除 key 和 value 两端的引号）
	return struct{ key, value string }{
		key:   strings.Trim(key.String(), `"`),
		value: strings.Trim(value.String(), `"`),
	}, nil
}

// SplitString splits a string along the specified separator, but it
// ignores anything between double quotes for splitting. We do simple
// inside/outside quote counting. Quotes are not stripped from output.
// 分割字符串，但尊重引号（不分割引号内的分隔符）
func SplitStringRespectingQuotes(s string, sep rune) []string {
	const escapeChar rune = '"'

	var parts []string
	var part string
	inQuotes := false

	for _, c := range s {
		if c == escapeChar {
			if inQuotes {
				inQuotes = false
			} else {
				inQuotes = true
			}
		}

		// 如果遇到分隔符 sep（如逗号 ,），且不在引号内
		// If we've gotten the separator rune, consider the previous part
		// complete, but only if we're outside of quoted sections
		if c == sep && !inQuotes {
			// 将当前子字符串加入结果
			parts = append(parts, part)
			// 重置当前子字符串
			part = ""
			continue
		}
		part += string(c)
	}

	// 将最后一个子字符串加入结果并返回
	return append(parts, part)
}

// ParseCommandLineList parses comma separated string lists which are passed
// in on the command line. Spaces are trimmed off both sides of result
// strings.
// 分割字符串，不考虑引号（用于普通列表解析），排除分隔后的空字符串
func ParseCommandLineList(input string) []string {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	splitInput := strings.Split(input, ",")
	args := make([]string, 0, len(splitInput))
	for _, s := range splitInput {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			args = append(args, s)
		}
	}
	return args
}
