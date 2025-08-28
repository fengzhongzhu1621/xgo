package regex

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestRegexMatchAllGroups 使用正则表达式匹配字符串中的所有子组并返回结果。
// func RegexMatchAllGroups(pattern, str string) [][]string
func TestRegexMatchAllGroups(t *testing.T) {
	pattern := `(\w+\.+\w+)@(\w+)\.(\w+)`
	str := "Emails: a.b@example.com and c.d@example.com"

	result := strutil.RegexMatchAllGroups(pattern, str)

	assert.Equal(t, []string{"a.b@example.com", "a.b", "example", "com"}, result[0], "result[0]")
	assert.Equal(t, []string{"c.d@example.com", "c.d", "example", "com"}, result[1], "result[1]")
}
