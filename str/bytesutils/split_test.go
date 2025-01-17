package bytesutils

import (
	"testing"

	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/testutil/assert"
)

func TestCut(t *testing.T) {
	// test for byteutil.Cut()
	b, a, ok := byteutil.Cut([]byte("age=123"), '=')
	assert.True(t, ok)
	assert.Eq(t, []byte("age"), b)
	assert.Eq(t, []byte("123"), a)

	b, a, ok = byteutil.Cut([]byte("age=123"), 'x')
	assert.False(t, ok)
	assert.Eq(t, []byte("age=123"), b)
	assert.Empty(t, a)

	// SafeCut
	b, a = byteutil.SafeCut([]byte("age=123"), '=')
	assert.Eq(t, []byte("age"), b)
	assert.Eq(t, []byte("123"), a)

	// SafeCuts
	b, a = byteutil.SafeCuts([]byte("age=123"), []byte{'='})
	assert.Eq(t, []byte("age"), b)
	assert.Eq(t, []byte("123"), a)
}
