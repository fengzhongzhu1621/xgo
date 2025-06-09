package stringutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gookit/goutil/fmtutil"
	"github.com/stretchr/testify/assert"
)

// TestConcat 从连接字符串。length是连接后字符串的长度。如果不确定，则传入0或负数。
// func Concat(length int, str ...string) string
func TestConcat(t *testing.T) {
	result1 := strutil.Concat(12, "Hello", " ", "World", "!")
	result2 := strutil.Concat(11, "Go", " ", "Language")
	result3 := strutil.Concat(0, "An apple a ", "day，", "keeps the", " doctor away")

	assert.Equal(t, "Hello World!", result1)
	assert.Equal(t, "Go Language", result2)
	assert.Equal(t, "An apple a day，keeps the doctor away", result3)
}

func TestArgsWithSpaces(t *testing.T) {
	args := []interface{}{"Hello", 123, true, 45.67}
	result := fmtutil.ArgsWithSpaces(args)
	assert.Equal(t, "Hello 123 true 45.67", result)
}
