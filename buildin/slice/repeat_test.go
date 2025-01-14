package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestRepeat Creates a slice with length n whose elements are passed param item.
// func Repeat[T any](item T, n int) []T
func TestRepeat(t *testing.T) {
	result := slice.Repeat("a", 3)

	assert.Equal(t, []string{"a", "a", "a"}, result)
}
