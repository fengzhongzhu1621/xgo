package hmac

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGenerateHmacSignedString(t *testing.T) {
	signed, _ := GenerateHmacSignedString("sha256", []byte("key"), []byte("1"))
	assert.Equal(t, signed, "bakfuRUXvh9c3POvkdfUDHF91jijBhV2BvsuWE966SY=")
}
