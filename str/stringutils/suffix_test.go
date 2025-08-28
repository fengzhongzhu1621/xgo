package stringutils

import (
	"strings"
	"testing"
)

func TestHasSuffix(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		suffix   string
		expected bool
	}{
		{
			name:     "后缀存在",
			s:        "Hello, World!",
			suffix:   "World!",
			expected: true,
		},
		{
			name:     "后缀不存在",
			s:        "Hello, World!",
			suffix:   "Universe!",
			expected: false,
		},
		{
			name:     "空字符串和空后缀",
			s:        "",
			suffix:   "",
			expected: true,
		},
		{
			name:     "s 为空，suffix 不为空",
			s:        "",
			suffix:   "Test",
			expected: false,
		},
		{
			name:     "suffix 为空",
			s:        "AnyString",
			suffix:   "",
			expected: true,
		},
		{
			name:     "s 和 suffix 都为空",
			s:        "",
			suffix:   "",
			expected: true,
		},
		{
			name:     "后缀为单个字符且存在",
			s:        "Golang is awesome",
			suffix:   "e",
			expected: true,
		},
		{
			name:     "后缀为单个字符但不存在",
			s:        "Python is great",
			suffix:   "e",
			expected: false,
		},
		{
			name:     "后缀与 s 完全相同",
			s:        "TestString",
			suffix:   "TestString",
			expected: true,
		},
		{
			name:     "后缀长度大于 s",
			s:        "Short",
			suffix:   "ShorterSuffix",
			expected: false,
		},
		{
			name:     "后缀在中间位置",
			s:        "This is a test string",
			suffix:   "test",
			expected: false,
		},
		{
			name:     "后缀在末尾但有其他字符",
			s:        "Hello, World!!",
			suffix:   "World!",
			expected: true,
		},
		{
			name:     "Unicode 字符后缀",
			s:        "こんにちは世界",
			suffix:   "世界",
			expected: true,
		},
		{
			name:     "Unicode 字符后缀不存在",
			s:        "こんにちは世界",
			suffix:   "世界！",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.HasSuffix(tt.s, tt.suffix)
			if result != tt.expected {
				t.Errorf(
					"HasSuffix(%q, %q) = %v; expected %v",
					tt.s,
					tt.suffix,
					result,
					tt.expected,
				)
			}
		})
	}
}
