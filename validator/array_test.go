package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateArray(t *testing.T) {
	t.Parallel()

	// not a slice
	a := "str"
	valid, message := ValidateArray(a)
	assert.False(t, valid)
	assert.NotEmpty(t, message)

	// a empty slice
	b := []string{}
	valid, message = ValidateArray(b)
	assert.False(t, valid)
	assert.Contains(t, message, "at least 1 item")

	// invalid
	type Test struct {
		Name string `json:"name" binding:"required"`
	}
	c := []Test{
		{""},
	}
	valid, message = ValidateArray(c)
	assert.False(t, valid)
	assert.Contains(t, message, "data in array")

	// valid
	d := []Test{
		{"aaaa"},
	}
	valid, message = ValidateArray(d)
	assert.True(t, valid)
	assert.Equal(t, "valid", message)
}
