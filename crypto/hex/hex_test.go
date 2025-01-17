package hex

import (
	"testing"

	"github.com/gookit/goutil/byteutil"
	"github.com/stretchr/testify/assert"
)

func TestHexEncoder(t *testing.T) {
	src := []byte("abc1234566")
	dst := byteutil.HexEncoder.Encode(src)
	assert.NotEmpty(t, dst)

	decSrc, err := byteutil.HexEncoder.Decode(dst)
	assert.NoError(t, err)
	assert.Equal(t, src, decSrc)
}
