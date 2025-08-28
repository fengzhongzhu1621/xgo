package assert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssert(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(123, 123, "they should be equal")

	assert.NotEqual(123, 456, "they should not be equal")

	assert.Nil(nil)

	assert.NotNil(1)
}
