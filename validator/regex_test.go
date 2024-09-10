package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidIDRegex(t *testing.T) {
	t.Parallel()

	assert.True(t, ValidIDRegex.MatchString("abc"))
	assert.True(t, ValidIDRegex.MatchString("abc-def"))
	assert.True(t, ValidIDRegex.MatchString("abc_def"))
	assert.True(t, ValidIDRegex.MatchString("abc_"))
	assert.True(t, ValidIDRegex.MatchString("abc-"))
	assert.True(t, ValidIDRegex.MatchString("abc-9"))
	assert.True(t, ValidIDRegex.MatchString("abc9ed"))

	assert.False(t, ValidIDRegex.MatchString("Abc"))
	assert.False(t, ValidIDRegex.MatchString("aBc"))
	assert.False(t, ValidIDRegex.MatchString("abC"))
	assert.False(t, ValidIDRegex.MatchString("_abc"))
	assert.False(t, ValidIDRegex.MatchString("*abc"))
	assert.False(t, ValidIDRegex.MatchString("9abc"))
	assert.False(t, ValidIDRegex.MatchString("abc+"))
	assert.False(t, ValidIDRegex.MatchString("abc*"))
	assert.False(t, ValidIDRegex.MatchString("42"))
}
