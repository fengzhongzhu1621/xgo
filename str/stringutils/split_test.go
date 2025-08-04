package stringutils

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHead 根据分隔符分割字符串.
func TestHead(t *testing.T) {
	s := "abc__def"
	sep := "__"
	left, right := Head(s, sep)
	assert.Equal(t, left, "abc")
	assert.Equal(t, right, "def")
}

// TestFieldsFunc 测试自定义分割
func TestFieldsFunc(t *testing.T) {
	s := "Hello, 世界! This is a test."

	// 用于将字符串按照指定的分隔函数进行分割，并返回一个字符串切片。这个函数非常灵活，因为它允许你自定义分隔符的判断逻辑。
	// 使用 FieldsFunc 分割字符串，以空白字符（空格、制表符等）和标点符号作为分隔符
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})

	fmt.Println(fields) // 输出：["Hello" "世界" "This" "is" "a" "test"]
}

// TestSplitEx 分割给定的字符串，无论结果是否包含空字符串。
// func SplitEx(s, sep string, removeEmptyString bool) []string
func TestSplitEx(t *testing.T) {
	type args struct {
		s                 string
		sep               string
		removeEmptyString bool
	}
	tests := []struct {
		name string
		args *args
		want []string
	}{
		{
			name: "test1",
			args: &args{
				s:                 " a b c ",
				sep:               "",
				removeEmptyString: true,
			},
			want: []string{},
		},
		{
			name: "test2",
			args: &args{
				s:                 " a b c ",
				sep:               " ",
				removeEmptyString: false,
			},
			want: []string{"", "a", "b", "c", ""},
		},
		{
			name: "test3",
			args: &args{
				s:                 " a b c ",
				sep:               " ",
				removeEmptyString: true,
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "test4",
			args: &args{
				s:                 "a = b = c = ",
				sep:               " = ",
				removeEmptyString: false,
			},
			want: []string{"a", "b", "c", ""},
		},
		{
			name: "test5",
			args: &args{
				s:                 "a = b = c = ",
				sep:               " = ",
				removeEmptyString: true,
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "test6",
			args: &args{
				s:                 "a = b = c = ",
				sep:               "=",
				removeEmptyString: false,
			},
			want: []string{"a ", " b ", " c ", " "},
		},
		{
			name: "test7",
			args: &args{
				s:                 "a = b = c =  ",
				sep:               "=",
				removeEmptyString: false,
			},
			want: []string{"a ", " b ", " c ", "  "},
		},
		{
			name: "test7",
			args: &args{
				s:                 "a = b = c =",
				sep:               "=",
				removeEmptyString: false,
			},
			want: []string{"a ", " b ", " c ", ""},
		},
		{
			name: "test9",
			args: &args{
				s:                 "a = b = c = ",
				sep:               "=",
				removeEmptyString: true,
			},
			want: []string{"a ", " b ", " c ", " "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.SplitEx(tt.args.s, tt.args.sep, tt.args.removeEmptyString)
			assert.Equal(t, tt.want, expect, tt.name)
		})
	}
}

// TestSplitWords 把一个字符串分割成若干单词，每个单词只包含字母字符。
// func SplitWords(s string) []string
func TestSplitWords(t *testing.T) {
	{
		type args struct {
			str string
		}
		tests := []struct {
			name string
			args args
			want []string
		}{
			{"", args{"PascalCase"}, []string{"Pascal", "Case"}},
			{"", args{"camelCase"}, []string{"camel", "Case"}},
			{"", args{"snake_case"}, []string{"snake", "case"}},
			{"", args{"kebab_case"}, []string{"kebab", "case"}},
			{"", args{"_test text_"}, []string{"test", "text"}},
			{"", args{"UPPERCASE"}, []string{"UPPERCASE"}},
			{"", args{"HTTPCode"}, []string{"HTTP", "Code"}},
			{"", args{"Int8Value"}, []string{"Int", "8", "Value"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equalf(t, tt.want, lo.Words(tt.args.str), "words(%v)", tt.args.str)
			})
		}
	}

	{
		type args struct {
			s string
		}
		tests := []struct {
			name string
			args *args
			want []string
		}{
			{
				name: "test1",
				args: &args{
					s: "a word",
				},
				want: []string{"a", "word"},
			},
			{
				name: "test2",
				args: &args{
					s: "I'am a programmer",
				},
				want: []string{"I'am", "a", "programmer"},
			},
			{
				name: "test3",
				args: &args{
					s: "Bonjour, je suis programmeur",
				},
				want: []string{"Bonjour", "je", "suis", "programmeur"},
			},
			{
				name: "test4",
				args: &args{
					s: "a -b-c' 'd'e",
				},
				want: []string{"a", "b-c'", "d'e"},
			},
			{
				name: "test5",
				args: &args{
					s: "你好，我是一名码农",
				},
				want: nil,
			},
			{
				name: "test6",
				args: &args{
					s: "こんにちは，私はプログラマーです",
				},
				want: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				expect := strutil.SplitWords(tt.args.s)
				assert.Equal(t, tt.want, expect, tt.name)
			})
		}
	}
}

// TestWordCount 返回有意义单词的数目，单词仅包含字母字符。
// func WordCount(s string) int
func TestWordCount(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args *args
		want int
	}{
		{
			name: "test1",
			args: &args{
				s: "a word",
			},
			want: 2,
		},
		{
			name: "test2",
			args: &args{
				s: "I'am a programmer",
			},
			want: 3,
		},
		{
			name: "test3",
			args: &args{
				s: "Bonjour, je suis programmeur",
			},
			want: 4,
		},
		{
			name: "test4",
			args: &args{
				s: "a -b-c' 'd'e",
			},
			want: 3,
		},
		{
			name: "test5",
			args: &args{
				s: "你好，我是一名码农",
			},
			want: 0,
		},
		{
			name: "test6",
			args: &args{
				s: "こんにちは，私はプログラマーです",
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.WordCount(tt.args.s)
			assert.Equal(t, tt.want, expect, tt.name)
		})
	}
}

// Lines 返回一个迭代器，该迭代器遍历字符串 s 中以换行符结束的行，生成的行包括它们的终止换行符。
// 如果 s 为空，迭代器将不产生任何行；
// 如果 s 末尾没有换行符，则最后生成的行将不以换行符结束。该迭代器为一次性使用。
func TestSplitLines(t *testing.T) {
	text := `第一行
第二行
第三行`

	lines := strings.Lines(text)

	// 打印每一行
	for line := range lines {
		fmt.Printf("%q\n", line)
	}

	// "第一行\n"
	// "第二行\n"
	// "第三行"
}

// SplitSeq 返回一个迭代器，遍历字符串 s 中由分隔符 sep 分隔的所有子字符串。
// 迭代器生成的字符串与使用 Split(s, sep) 返回的字符串相同，但不构造切片。该迭代器为一次性使用。
func TestSplitSeq(t *testing.T) {
	// 使用特定分隔符分割字符串
	s := "a,b,c,d"
	fmt.Println("Split string by comma:")
	for part := range strings.SplitSeq(s, ",") {
		fmt.Printf("%q\n", part)
	}
	// "a"
	// "b"
	// "c"
	// "d"

	// 使用空分隔符分割成字符
	text := "Hello世界"
	fmt.Println("\nSplit into characters:")
	for char := range strings.SplitSeq(text, "") {
		fmt.Printf("%q\n", char)
	}
	// "H"
	// "e"
	// "l"
	// "l"
	// "o"
	// "世"
	// "界"
}

func TestSplitAfterSeq(t *testing.T) {
	// 使用分隔符分割(并会保留分隔符)
	s := "a,b,c,d"
	fmt.Println("Split string by comma (keeping separators):")
	for part := range strings.SplitAfterSeq(s, ",") {
		fmt.Printf("%q\n", part)
	}
	// "a,"
	// "b,"
	// "c,"
	// "d"
}

// SplitSeq 和 SplitAfterSeq 使用指定的分隔符(separator)来分割字符串
// FieldsSeq 自动使用空白字符(whitespace)作为分隔符，包括空格、制表符、换行符等
//
// SplitSeq 和 SplitAfterSeq 会保留空字符串(在连续分隔符之间)
// FieldsSeq 会忽略连续的空白字符，不会产生空字符串
// SplitSeq "a  b" 使用 " " 分割会产生: ["a", "", "b"]
// FieldsSeq "a  b" 会产生: ["a", "b"]
func TestFieldsSeq(t *testing.T) {
	// 通过空格分割
	text := "a b c d"
	fmt.Println("Split string into fields:")
	for word := range strings.FieldsSeq(text) {
		fmt.Printf("%q\n", word)
	}
	// "a"
	// "b"
	// "c"
	// "d"

	// 通过多个空格来分割
	textWithSpaces := "  a   b   c  "
	fmt.Println("\nSplit string with multiple spaces:")
	for word := range strings.FieldsSeq(textWithSpaces) {
		fmt.Printf("%q\n", word)
	}
	// "a"
	// "b"
	// "c"
}

func TestFieldsFuncSeq(t *testing.T) {
	// 使用空格分割 (和FieldsSeq效果类似)
	text := "a b c  d"
	fmt.Println("Split on whitespace:")
	for word := range strings.FieldsFuncSeq(text, unicode.IsSpace) {
		fmt.Printf("%q\n", word)
	}
	// "a"
	// "b"
	// "c"
	// "d"

	// 根据数字切割
	mixedText := "abc123def456ghi"
	fmt.Println("\nSplit on digits:")
	for word := range strings.FieldsFuncSeq(mixedText, unicode.IsDigit) {
		fmt.Printf("%q\n", word)
	}
	// "abc"
	// "def"
	// "ghi"
}

func TestParseInputMapping(t *testing.T) {
	var src string
	var expected, parsed map[string]string
	var err error

	src = "key1:value1,key2:value2"
	expected = map[string]string{"key1": "value1", "key2": "value2"}
	parsed, err = ParseCommandlineMap(src)
	require.NoError(t, err)
	assert.Equal(t, expected, parsed)

	src = `key1:"value1,value2",key2:value3`
	expected = map[string]string{"key1": "value1,value2", "key2": "value3"}
	parsed, err = ParseCommandlineMap(src)
	require.NoError(t, err)
	assert.Equal(t, expected, parsed)

	src = `key1:"value1,value2,key2:value3"`
	expected = map[string]string{"key1": "value1,value2,key2:value3"}
	parsed, err = ParseCommandlineMap(src)
	require.NoError(t, err)
	assert.Equal(t, expected, parsed)

	src = `"key1,key2":value1`
	expected = map[string]string{"key1,key2": "value1"}
	parsed, err = ParseCommandlineMap(src)
	require.NoError(t, err)
	assert.Equal(t, expected, parsed)
}

func TestSplitStringRespectingQuotes(t *testing.T) {
	var src string
	var expected, result []string
	var err error

	src = "1,2,3"
	expected = []string{"1", "2", "3"}
	result = SplitStringRespectingQuotes(src, ',')
	require.NoError(t, err)
	assert.Equal(t, expected, result)

	src = `"1,2",3`
	expected = []string{`"1,2"`, "3"}
	result = SplitStringRespectingQuotes(src, ',')
	require.NoError(t, err)
	assert.Equal(t, expected, result)

	src = `1,"2,3",`
	expected = []string{"1", `"2,3"`, ""}
	result = SplitStringRespectingQuotes(src, ',')
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
