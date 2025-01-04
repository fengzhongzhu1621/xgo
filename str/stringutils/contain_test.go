package stringutils

import (
	"testing"
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
	}{
		{"hello", "lo", true},
		{"abcdefg", "cf", true},
		{"abcdefg", "a", true},
		{"abcdefg", "b", true},
		{"abcdefg", "cfa", false},
		{"abcdefg", "aa", false},
		{"世界", "a", false},
		{"Hello 世界", "界", true},
		{"Hello 世界", "elo", true},
	}
	for _, v := range tests {
		res := SublimeContains(v.text, v.substr)
		if res != v.pass {
			t.Fatalf("Failed: %v - res:%v", v, res)
		}
	}
}
