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
