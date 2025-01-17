package validator

import (
	"testing"

	"github.com/gookit/goutil/byteutil"

	"github.com/stretchr/testify/assert"
)

func TestHasPrefixAndSuffix(t *testing.T) {
	s := []byte("ab____cd")
	prefix := []byte("ab")
	suffix := []byte("cd")
	actual := HasPrefixAndSuffix(s, prefix, suffix)
	assert.Equal(t, actual, true, "they should be equal")

	prefix = []byte("abc")
	actual = HasPrefixAndSuffix(s, prefix, suffix)
	assert.Equal(t, actual, false, "they should not be equal")
}

func TestIsNumChar(t *testing.T) {
	tests := []struct {
		args byte
		want bool
	}{
		{'2', true},
		{'a', false},
		{'+', false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, byteutil.IsNumChar(tt.args))
	}
}
