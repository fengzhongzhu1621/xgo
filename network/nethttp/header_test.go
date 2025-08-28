package nethttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendEnv(t *testing.T) {
	env := []string{"a", "b"}
	key := "key"
	value1 := "value1"
	value2 := "value2"
	env1 := AppendEnv(env, key, value1, value2)
	assert.Equal(t, []string{"a", "b"}, env)
	assert.Equal(t, []string{"a", "b", "KEY=value1, value2"}, env1)
}

var mimetest = [][3]string{
	{"Content-Type: text/plain", "Content-Type", "text/plain"},
	{"Content-Type:    ", "Content-Type", ""},
}

func TestSplitMimeHeader(t *testing.T) {
	for _, tst := range mimetest {
		s, v := SplitMimeHeader(tst[0])
		if tst[1] != s || tst[2] != v {
			t.Errorf("%v and %v  are not same as expected %v and %v", s, v, tst[1], tst[2])
		}
	}
}
