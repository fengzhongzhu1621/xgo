package stringutils

import (
	"strings"
	"testing"

	"github.com/samber/lo"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// BenchmarkToLower 性能较好
func BenchmarkToLower(b *testing.B) {
	string1 := "ABCDEFGHIJKLMN"
	string2 := "abcdefghijklmn"
	for i := 0; i < b.N; i++ {
		ToLower(string1)
		ToLower(string2)
	}
}

func BenchmarkBuildInToLower(b *testing.B) {
	string1 := "ABCDEFGHIJKLMN"
	string2 := "abcdefghijklmn"
	for i := 0; i < b.N; i++ {
		strings.ToLower(string1)
		strings.ToLower(string2)
	}
}

// TestCapitalize 将字符串的第一个字符转换为大写
func TestCapitalize(t *testing.T) {
	{
		type args struct {
			word string
		}
		tests := []struct {
			name string
			args args
			want string
		}{
			{"", args{"hello"}, "Hello"},
			{"", args{"heLLO"}, "Hello"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equalf(t, tt.want, lo.Capitalize(tt.args.word), "Capitalize(%v)", tt.args.word)
			})
		}
	}

	{
		tests := []struct {
			name string
			args string
			want string
		}{
			{
				name: "test1",
				args: "",
				want: "",
			},
			{
				name: "test2",
				args: "Foo",
				want: "Foo",
			},
			{
				name: "test3",
				args: "_foo",
				want: "_foo",
			},
			{
				name: "test4",
				args: "fooBar",
				want: "Foobar",
			},
			{
				name: "test5",
				args: "foo-bar",
				want: "Foo-bar",
			},
			{
				name: "test6",
				args: "convert_test.go",
				want: "Convert_test.go",
			},
			{
				name: "test7",
				args: "convertTest.go",
				want: "Converttest.go",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				expect := strutil.Capitalize(tt.args)
				assert.Equal(t, tt.want, expect)
			})
		}
	}
}

// TestLowerFirst 将字符串的第一个字符转换为小写。
func TestLowerFirst(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "test1",
			args: "",
			want: "",
		},
		{
			name: "test2",
			args: "bar",
			want: "bar",
		},
		{
			name: "test3",
			args: "_foo",
			want: "_foo",
		},
		{
			name: "test4",
			args: "BAr",
			want: "bAr",
		},
		{
			name: "test5",
			args: "Bar大",
			want: "bar大",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.LowerFirst(tt.args)
			assert.Equal(t, tt.want, expect)
		})
	}
}

// TestUpperFirst 将字符串的第一个字符转换为小写。
func TestUpperFirst(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "test1",
			args: "",
			want: "",
		},
		{
			name: "test2",
			args: "bar",
			want: "Bar",
		},
		{
			name: "test3",
			args: "_foo",
			want: "_foo",
		},
		{
			name: "test4",
			args: "BAr",
			want: "BAr",
		},
		{
			name: "test5",
			args: "bar大",
			want: "Bar大",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.UpperFirst(tt.args)
			assert.Equal(t, tt.want, expect)
		})
	}
}
