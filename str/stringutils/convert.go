package stringutils

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/fengzhongzhu1621/xgo/validator"
)

var (
	// 这是一个指向正则表达式对象的指针，用于匹配路径参数。在后续的代码中，这个正则表达式会被用来从 URL 路径中提取参数。
	pathParamRE *regexp.Regexp
	// 这是一个映射（map），用于存储 Go 语言的预声明标识符。预声明标识符是 Go 语言中预先定义的一些关键字和基本类型，比如 int、string、true 等。
	predeclaredSet map[string]struct{}
	// 用于存储分隔符字符。这些字符在处理字符串时可能会被用作分隔符。
	separatorSet map[rune]struct{}
)

func init() {
	pathParamRE = regexp.MustCompile(`{[.;?]?([^{}*]+)\*?}`)

	predeclaredIdentifiers := []string{
		// Types
		"bool",
		"byte",
		"complex64",
		"complex128",
		"error",
		"float32",
		"float64",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"rune",
		"string",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		// Constants
		"true",
		"false",
		"iota",
		// Zero value
		"nil",
		// Functions
		"append",
		"cap",
		"close",
		"complex",
		"copy",
		"delete",
		"imag",
		"len",
		"make",
		"new",
		"panic",
		"print",
		"println",
		"real",
		"recover",
	}

	// 预声明标识符
	predeclaredSet = map[string]struct{}{}
	for _, id := range predeclaredIdentifiers {
		predeclaredSet[id] = struct{}{}
	}

	// 分隔符字符
	separators := "-#@!$&=.+:;_~ (){}[]"
	separatorSet = map[rune]struct{}{}
	for _, r := range separators {
		separatorSet[r] = struct{}{}
	}
}

// ToLower 字符串转换为小写，在转化前先判断是否包含大写字符，比strings.ToLower性能高.
func ToLower(s string) string {
	// 判断字符串是否包含小写字母
	if validator.IsLower(s) {
		return s
	}
	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	// []bytes转换为字符串
	return cast.BytesToString(b)
}

// UnicodeTitle 首字母大写.
func UnicodeTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToUpper(v)) + s[k+1:]
	}
	return ""
}

// UnicodeUnTitle 首字母小写.
func UnicodeUnTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToLower(v)) + s[k+1:]
	}
	return ""
}

// UppercaseFirstCharacter Uppercases the first character in a string. This assumes UTF-8, so we have
// to be careful with unicode, don't treat it as a byte array.
// 将字符串的第一个字符转换为大写，适用于需要首字母大写的场景，如将 http 转换为 Http。
func UppercaseFirstCharacter(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Uppercase the first character in a identifier with pkg name. This assumes UTF-8, so we have
// to be careful with unicode, don't treat it as a byte array.
// 处理带有包名的字符串，仅将最后一个部分的第一个字符大写，如将 pkg.http 转换为 pkg.Http。
func UppercaseFirstCharacterWithPkgName(str string) string {
	if str == "" {
		return ""
	}

	segs := strings.Split(str, ".")
	var prefix string
	if len(segs) == 2 {
		prefix = segs[0] + "."
		str = segs[1]
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return prefix + string(runes)
}

// LowercaseFirstCharacter Lowercases the first character in a string. This assumes UTF-8, so we have
// to be careful with unicode, don't treat it as a byte array.
// 将字符串的第一个字符转换为小写
func LowercaseFirstCharacter(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// Lowercase the first upper characters in a string for case of abbreviation.
// This assumes UTF-8, so we have to be careful with unicode, don't treat it as a byte array.
// 将字符串中连续的大写字母序列的第一个字母转换为小写，适用于缩写处理，如将 HTTPResponse 转换为 HttpResponse。
func LowercaseFirstCharacters(str string) string {
	if str == "" {
		return ""
	}

	runes := []rune(str)

	for i := 0; i < len(runes); i++ {
		next := i + 1
		if i != 0 && next < len(runes) && unicode.IsLower(runes[next]) {
			break
		}

		runes[i] = unicode.ToLower(runes[i])
	}

	return string(runes)
}

// ChangeInitialCase 将字符串的首字符转换为指定格式
func ChangeInitialCase(s string, mapper func(rune) rune) string {
	if s == "" {
		return s
	}
	// 返回第一个utf8字符，n是字符的长度，即返回首字符
	r, n := utf8.DecodeRuneInString(s)
	// 根据mapper方法转换首字符
	return string(mapper(r)) + s[n:]
}

// Normalize to trim space of the str and get it's upper format
// for example, Normalize(" hello world") ==> "HELLO WORLD"
func Normalize(str string) string {
	return strings.ToUpper(strings.TrimSpace(str))
}

// ToCamelCase will convert query-arg style strings to CamelCase. We will
// use `., -, +, :, ;, _, ~, ' ', (, ), {, }, [, ]` as valid delimiters for words.
// So, "word.word-word+word:word;word_word~word word(word)word{word}[word]"
// would be converted to WordWordWordWordWordWordWordWordWordWordWordWordWord
// 将特定格式的字符串转换为驼峰式（CamelCase）命名
// 将类似查询参数风格的字符串转换为驼峰式（CamelCase）。
// 有效的单词分隔符包括 ., -, +, :, ;, _, ~, ' ', (, ), {, }, [, ]。
// 因此，输入字符串中的这些分隔符将被移除，并将分隔后的单词首字母大写以形成驼峰式命名。
func ToCamelCase(str string) string {
	s := strings.Trim(str, " ")

	n := ""
	capNext := true
	for _, v := range s {
		if unicode.IsUpper(v) {
			n += string(v)
		}
		if unicode.IsDigit(v) {
			n += string(v)
		}
		if unicode.IsLower(v) {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}

		// 检查当前字符 v 是否是分隔符
		// 如果是分隔符，capNext 被设置为 true，表示下一个非分隔符字符需要大写。
		// 如果不是分隔符，capNext 保持不变。
		_, capNext = separatorSet[v]
	}
	return n
}

// ToCamelCaseWithDigits function will convert query-arg style strings to CamelCase. We will
// use `., -, +, :, ;, _, ~, ' ', (, ), {, }, [, ]` as valid delimiters for words.
// The difference of ToCamelCase that letter after a number becomes capitalized.
// So, "word.word-word+word:word;word_word~word word(word)word{word}[word]3word"
// would be converted to WordWordWordWordWordWordWordWordWordWordWordWordWord3Word
// 与 ToCamelCase 类似，将查询参数风格的字符串转换为驼峰式（CamelCase），
// 但有一个关键区别：在数字之后的字母会被大写。这意味着数字被视为单词的一部分，并且其后紧跟的字母应大写以开始新的单词部分。
func ToCamelCaseWithDigits(s string) string {
	res := bytes.NewBuffer(nil)
	capNext := true
	for _, v := range s {
		if unicode.IsUpper(v) {
			res.WriteRune(v)
			capNext = false
			continue
		}
		if unicode.IsDigit(v) {
			res.WriteRune(v)
			capNext = true
			continue
		}
		if unicode.IsLower(v) {
			if capNext {
				res.WriteRune(unicode.ToUpper(v))
			} else {
				res.WriteRune(v)
			}
			capNext = false
			continue
		}
		capNext = true
	}
	return res.String()
}
