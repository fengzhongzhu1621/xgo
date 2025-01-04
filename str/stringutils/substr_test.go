package stringutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestSubstring 返回从指定偏移位置开始的指定长度的子字符串。
// func Substring(s string, offset int, length uint) string
func TestSubstring(t *testing.T) {
	type args struct {
		s      string
		offset int
		length uint
	}
	tests := []struct {
		name string
		args *args
		want string
	}{
		{
			name: "test1",
			args: &args{
				s:      "abcde",
				offset: 1,
				length: 3,
			},
			want: "bcd",
		},
		{
			name: "test2",
			args: &args{
				s:      "abcde",
				offset: 1,
				length: 5,
			},
			want: "bcde",
		},
		{
			name: "test3",
			args: &args{
				s:      "abcde",
				offset: -1,
				length: 3,
			},
			want: "e",
		},
		{
			name: "test4",
			args: &args{
				s:      "abcde",
				offset: -2,
				length: 2,
			},
			want: "de",
		},
		{
			name: "test5",
			args: &args{
				s:      "abcde",
				offset: -2,
				length: 3,
			},
			want: "de",
		},
		{
			name: "test6",
			args: &args{
				s:      "你好，欢迎你",
				offset: 0,
				length: 2,
			},
			want: "你好",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.Substring(tt.args.s, tt.args.offset, tt.args.length)
			assert.Equal(t, tt.want, expect, tt.name)
		})
	}
}

// TestIndexOffset 在字符串中从 idxFrom 偏移的位置开始，返回子字符串 substr 的第一个实例的索引；如果 substr 不存在于字符串中，则返回 -1。
// func IndexOffset(str string, substr string, idxFrom int) int
func TestIndexOffset(t *testing.T) {
	str := "foo bar hello world"

	result1 := strutil.IndexOffset(str, "o", 5)
	result2 := strutil.IndexOffset(str, "o", 0)
	result3 := strutil.IndexOffset(str, "d", len(str)-1)
	result4 := strutil.IndexOffset(str, "d", len(str))
	result5 := strutil.IndexOffset(str, "f", -1)

	assert.Equal(t, 12, result1)
	assert.Equal(t, 1, result2)
	assert.Equal(t, 18, result3)
	assert.Equal(t, -1, result4)
	assert.Equal(t, -1, result5)
}

// TestSubInBetween 返回源字符串中开始和结束位置（不包括）之间的子字符串。
// func SubInBetween(str string, start string, end string) string
func TestSubInBetween(t *testing.T) {
	str := "abcded 你好啊"

	result1 := strutil.SubInBetween(str, "", "de")
	result2 := strutil.SubInBetween(str, "a", "d")
	result3 := strutil.SubInBetween(str, "你", "啊")

	assert.Equal(t, "abc", result1)
	assert.Equal(t, "bc", result2)
	assert.Equal(t, "好", result3)
}
