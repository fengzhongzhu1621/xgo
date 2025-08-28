package bytesutils

import (
	"testing"

	"github.com/gookit/goutil/byteutil"
	"github.com/stretchr/testify/assert"
)

func TestCut(t *testing.T) {
	// test for byteutil.Cut()
	b, a, ok := byteutil.Cut([]byte("age=123"), '=')
	assert.True(t, ok)
	assert.Equal(t, []byte("age"), b)
	assert.Equal(t, []byte("123"), a)

	b, a, ok = byteutil.Cut([]byte("age=123"), 'x')
	assert.False(t, ok)
	assert.Equal(t, []byte("age=123"), b)
	assert.Empty(t, a)

	// SafeCut
	b, a = byteutil.SafeCut([]byte("age=123"), '=')
	assert.Equal(t, []byte("age"), b)
	assert.Equal(t, []byte("123"), a)

	// SafeCuts
	b, a = byteutil.SafeCuts([]byte("age=123"), []byte{'='})
	assert.Equal(t, []byte("age"), b)
	assert.Equal(t, []byte("123"), a)
}
