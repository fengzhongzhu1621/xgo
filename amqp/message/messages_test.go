package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessages_IDs(t *testing.T) {
	msgs := Messages{
		NewMessage("1", nil),
		NewMessage("2", nil),
		NewMessage("3", nil),
	}

	assert.Equal(t, []string{"1", "2", "3"}, msgs.IDs())
}
