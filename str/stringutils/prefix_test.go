package stringutils

import (
	"strings"
	"testing"
)

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		prefix   string
		expected bool
	}{
		{
			name:     "前缀存在",
			s:        "Hello, World!",
			prefix:   "Hello",
			expected: true,
		},
		{
			name:     "前缀不存在",
			s:        "Hello, World!",
			prefix:   "Hi",
			expected: false,
		},
		{
			name:     "空字符串和空前缀",
			s:        "",
			prefix:   "",
			expected: true,
		},
		{
			name:     "s 为空，prefix 不为空",
			s:        "",
			prefix:   "Test",
			expected: false,
		},
		{
			name:     "prefix 为空",
			s:        "AnyString",
			prefix:   "",
			expected: true,
		},
		{
			name:     "s 和 prefix 都为空",
			s:        "",
			prefix:   "",
			expected: true,
		},
		{
			name:     "前缀为单个字符且存在",
			s:        "Golang is awesome",
			prefix:   "G",
			expected: true,
		},
		{
			name:     "前缀为单个字符但不存在",
			s:        "Python is great",
			prefix:   "G",
			expected: false,
		},
		{
			name:     "前缀与 s 完全相同",
			s:        "TestString",
			prefix:   "TestString",
			expected: true,
		},
		{
			name:     "前缀长度大于 s",
			s:        "Short",
			prefix:   "ShorterPrefix",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.HasPrefix(tt.s, tt.prefix)
			if result != tt.expected {
				t.Errorf(
					"HasPrefix(%q, %q) = %v; expected %v",
					tt.s,
					tt.prefix,
					result,
					tt.expected,
				)
			}
		})
	}
}
