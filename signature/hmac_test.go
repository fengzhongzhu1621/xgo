package signature

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGenereateHmacSignedString(t *testing.T) {
	signed, _ := GenereateHmacSignedString("sha256", []byte("key"), []byte("1"))
	assert.Equal(t, signed, "bakfuRUXvh9c3POvkdfUDHF91jijBhV2BvsuWE966SY=")
}

