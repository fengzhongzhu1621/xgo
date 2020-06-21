package bytes_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
