package bytesutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateBytes(t *testing.T) {
	content := []byte("Hello, world!")
	truncated := TruncateBytes(content, 5)
	assert.Equal(t, "Hello", string(truncated))
}
