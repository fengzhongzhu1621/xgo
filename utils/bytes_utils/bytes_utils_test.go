package bytes_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var eol_tests = []string{
	"", "\n", "\r\n", "ok\n", "ok\n",
	"quite long string for our test\n",
	"quite long string for our test\r\n",
}

var eol_answers = []string{
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


func TestTrimEOL(t *testing.T) {
	for n := 0; n < len(eol_tests); n++ {
		answ := TrimEOL([]byte(eol_tests[n]))
		if string(answ) != eol_answers[n] {
			t.Errorf("Answer '%s' did not match predicted '%s'", answ, eol_answers[n])
		}
	}
}


func BenchmarkTrimEOL(b *testing.B) {
	for n := 0; n < b.N; n++ {
		TrimEOL([]byte(eol_tests[n%len(eol_tests)]))
	}
}
