package stringutils

import (
	"strings"
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestStringsMatch 判断 key 是否在 matchList 中
func TestStringsMatch(t *testing.T) {
	// Test case 1: key is empty
	if !StringsMatch("", "hello", "world") {
		t.Errorf("Expected true but got false")
	}
	// Test case 2: key is present in matchList
	if !StringsMatch("world", "hello", "world") {
		t.Errorf("Expected true but got false")
	}
	// Test case 3: key is not present in matchList
	if StringsMatch("foo", "hello", "world") {
		t.Errorf("Expected false but got true")
	}
	// Test case 4: matchList is empty
	if StringsMatch("hello") {
		t.Errorf("Expected false but got true")
	}
}

// TestStringsMatchObscure 检查 key 是否包含在由 matchList 中的字符串拼接而成的字符串中，且在比较时不区分大小写。
func TestStringsMatchObscure(t *testing.T) {
	// Test case 1: key is empty
	if !StringsMatchObscure("", "hello", "world") {
		t.Errorf("Expected true but got false")
	}
	// Test case 2: key is present in matchList
	if !StringsMatchObscure("world", "Hello", "WORLD") {
		t.Errorf("Expected true but got false")
	}
	// Test case 3: key is not present in matchList
	if StringsMatchObscure("foo", "hello", "world") {
		t.Errorf("Expected false but got true")
	}
	// Test case 4: matchList is empty
	if StringsMatchObscure("hello") {
		t.Errorf("Expected false but got true")
	}
}

func TestSublimeContains(t *testing.T) {
	tests := []struct {
		text   string
		substr string
		pass   bool
		pass1  bool
	}{
		{"hello", "lo", true, true},
		{"abcdefg", "cf", true, false},
		{"abcdefg", "a", true, true},
		{"abcdefg", "b", true, true},
		{"abcdefg", "cfa", false, false},
		{"abcdefg", "aa", false, false},
		{"世界", "a", false, false},
		{"Hello 世界", "界", true, true},
		{"Hello 世界", "elo", true, false},
	}
	for _, v := range tests {
		res := SublimeContains(v.text, v.substr)
		if res != v.pass {
			t.Fatalf("SublimeContains Failed: %v - res:%v", v, res)
		}
		res1 := strings.Contains(v.text, v.substr)
		if res1 != v.pass1 {
			t.Fatalf("Contains Failed: %v - res:%v", v, res1)
		}
	}
}

// TestContainsAll 如果目标字符串包含所有的子字符串，则返回true。
// func ContainsAll(str string, substrs []string) bool
func TestContainsAll(t *testing.T) {
	str := "hello world"

	result1 := strutil.ContainsAll(str, []string{"hello", "world"})
	result2 := strutil.ContainsAll(str, []string{"hello", "abc"})

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}

// TestContainsAny 如果目标字符串包含其中任意一个子字符串，则返回true。
// func ContainsAny(str string, substrs []string) bool
func TestContainsAny(t *testing.T) {
	str := "hello world"

	result1 := strutil.ContainsAny(str, []string{"hello", "world"})
	result2 := strutil.ContainsAny(str, []string{"hello", "abc"})
	result3 := strutil.ContainsAny(str, []string{"123", "abc"})

	assert.Equal(t, true, result1)
	assert.Equal(t, true, result2)
	assert.Equal(t, false, result3)
}
