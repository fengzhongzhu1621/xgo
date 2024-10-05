package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateBytes(t *testing.T) {
	content := []byte("Hello, world!")
	truncated := TruncateBytes(content, 5)
	assert.Equal(t, "Hello", string(truncated))
}

func TestTruncateBytesToString(t *testing.T) {
	content := []byte("Hello, world!")
	truncatedStr := TruncateBytesToString(content, 5)
	assert.Equal(t, "Hello", string(truncatedStr))
}
