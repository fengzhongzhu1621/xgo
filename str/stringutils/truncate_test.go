package stringutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestAfter 测试返回第一个匹配到的字符串后面的子串
func TestAfter(t *testing.T) {
	result1 := strutil.After("foo", "")
	result2 := strutil.After("foo", "foo")
	result3 := strutil.After("foo/bar", "foo")
	result4 := strutil.After("foo/bar", "/")
	result5 := strutil.After("foo/bar/baz", "/")
	result6 := strutil.After("/fo/foo/bar/foo/baz", "foo")

	assert.Equal(t, "foo", result1)
	assert.Equal(t, "", result2)
	assert.Equal(t, "/bar", result3)
	assert.Equal(t, "bar", result4)
	assert.Equal(t, "bar/baz", result5)
	assert.Equal(t, "/bar/foo/baz", result6)
}

// TestAfterLast 测试返回最后一个匹配到的字符串后面的子串
func TestAfterLast(t *testing.T) {
	result1 := strutil.AfterLast("foo", "")
	result2 := strutil.AfterLast("foo", "foo")
	result3 := strutil.AfterLast("foo/bar", "/")
	result4 := strutil.AfterLast("foo/bar/baz", "/")
	result5 := strutil.AfterLast("foo/bar/foo/baz", "foo")

	assert.Equal(t, "foo", result1)
	assert.Equal(t, "", result2)
	assert.Equal(t, "bar", result3)
	assert.Equal(t, "baz", result4)
	assert.Equal(t, "/baz", result5)
}

func TestBefore(t *testing.T) {
	result1 := strutil.Before("foo", "")
	result2 := strutil.Before("foo", "foo")
	result3 := strutil.Before("foo/bar", "/")
	result4 := strutil.Before("foo/bar/baz", "/")

	assert.Equal(t, "foo", result1)
	assert.Equal(t, "", result2)
	assert.Equal(t, "foo", result3)
	assert.Equal(t, "foo", result4)
}

func TestBeforeAll(t *testing.T) {
	result1 := strutil.BeforeLast("foo", "")
	result2 := strutil.BeforeLast("foo", "foo")
	result3 := strutil.BeforeLast("foo/bar", "/")
	result4 := strutil.BeforeLast("foo/bar/baz", "/")

	assert.Equal(t, "foo", result1)
	assert.Equal(t, "", result2)
	assert.Equal(t, "foo", result3)
	assert.Equal(t, "foo/bar", result4)
}

// TestWrap 用给定的字符串包裹一个字符串。
// func Wrap(str string, wrapWith string) string
func TestWrap(t *testing.T) {
	type args struct {
		s        string
		wrapWith string
	}
	tests := []struct {
		name string
		args *args
		want string
	}{
		{
			name: "test1",
			args: &args{
				s:        "foo",
				wrapWith: "",
			},
			want: "foo",
		},
		{
			name: "test2",
			args: &args{
				s:        "foo",
				wrapWith: "*",
			},
			want: "*foo*",
		},
		{
			name: "test3",
			args: &args{
				s:        "'foo'",
				wrapWith: "'",
			},
			want: "''foo''",
		},
		{
			name: "test4",
			args: &args{
				s:        "",
				wrapWith: "*",
			},
			want: "",
		},
		{
			name: "test5",
			args: &args{
				s:        "foo",
				wrapWith: "<>",
			},
			want: "<>foo<>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.Wrap(tt.args.s, tt.args.wrapWith)
			assert.Equal(t, tt.want, expect, tt.name)
		})
	}
}

// TestUnwrap 从另一字符串中解包一个给定的字符串。将会更改源字符串。
// func Unwrap(str string, wrapToken string) string
func TestUnwrap(t *testing.T) {
	type args struct {
		s         string
		wrapToken string
	}
	tests := []struct {
		name string
		args *args
		want string
	}{
		{
			name: "test1",
			args: &args{
				s:         "foo",
				wrapToken: "",
			},
			want: "foo",
		},
		{
			name: "test2",
			args: &args{
				s:         "*foo*",
				wrapToken: "*",
			},
			want: "foo",
		},
		{
			name: "test3",
			args: &args{
				s:         "*foo",
				wrapToken: "*",
			},
			want: "*foo",
		},
		{
			name: "test4",
			args: &args{
				s:         "foo*",
				wrapToken: "*",
			},
			want: "foo*",
		},
		{
			name: "test5",
			args: &args{
				s:         "**foo**",
				wrapToken: "*",
			},
			want: "*foo*",
		},
		{
			name: "test8",
			args: &args{
				s:         "**foo**",
				wrapToken: "**",
			},
			want: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.Unwrap(tt.args.s, tt.args.wrapToken)
			assert.Equal(t, tt.want, expect, tt.name)
		})
	}
}

// TestRemoveNonPrintable 从字符串中移除不可打印的字符。
// func RemoveNonPrintable(str string) string
func TestRemoveNonPrintable(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args *args
		want string
	}{
		{
			name: "test1",
			args: &args{
				s: "hello\u00a0 \u200bworld\n",
			},
			want: "hello world",
		},
		{
			name: "test2",
			args: &args{
				s: "你好😄",
			},
			want: "你好😄",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.RemoveNonPrintable(tt.args.s)
			assert.Equal(t, tt.want, expect, tt.name)
		})
	}
}
