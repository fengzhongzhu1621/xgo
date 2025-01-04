package stringutils

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

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
