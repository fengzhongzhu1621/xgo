package bytesutils

import (
	"testing"

	"github.com/gookit/goutil/byteutil"
	"github.com/stretchr/testify/assert"
)

func TestFirstLine(t *testing.T) {
	bs := []byte("hi\ninhere")
	assert.Equal(t, []byte("hi"), byteutil.FirstLine(bs))
	assert.Equal(t, []byte("hi"), byteutil.FirstLine([]byte("hi")))
}
