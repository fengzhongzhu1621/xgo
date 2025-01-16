package bytesutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var eolTests = []string{
	"", "\n", "\r\n", "ok\n", "ok\n",
	"quite long string for our test\n",
	"quite long string for our test\r\n",
}

var eolAnswers = []string{
	"", "", "", "ok", "ok",
	"quite long string for our test", "quite long string for our test",
}

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
