package validator

import (
	"testing"

	"github.com/duke-git/lancet/v2/validator"
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

// TestIsRegexMatch 检查字符串是否与正则表达式匹配。
// func IsRegexMatch(s, regex string) bool
func TestIsRegexMatch(t *testing.T) {
	result1 := validator.IsRegexMatch("abc", `^[a-zA-Z]+$`)
	result2 := validator.IsRegexMatch("ab1", `^[a-zA-Z]+$`)

	assert.Equal(t, true, result1)
	assert.Equal(t, false, result2)
}
