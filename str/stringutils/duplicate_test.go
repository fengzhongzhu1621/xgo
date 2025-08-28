package stringutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDuplicateElement(t *testing.T) {
	items := []string{"a", "b", "a"}
	dropDuplicatedItems := RemoveDuplicateElement(items)
	assert.Equal(t, dropDuplicatedItems, []string{"a", "b"})
}
