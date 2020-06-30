package string_utils

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestHead(t *testing.T) {
	s := "abc__def"
	sep := "__"
	left, right := Head(s, sep)
	assert.Equal(t, left, "abc")
	assert.Equal(t, right, "def")
}

func TestRemoveDuplicateElement(t *testing.T) {
	items := []string{"a", "b", "a"}
	dropDuplicatedItems := RemoveDuplicateElement(items)
	assert.Equal(t, dropDuplicatedItems, []string{"a", "b"})
}
