package string_utils

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestHead(t *testing.T) {
	s := "abc__def"
	sep := "__"
	left, right := head(s, sep)
	assert.Equal(t, left, "abc")
	assert.Equal(t, right, "def")
}

