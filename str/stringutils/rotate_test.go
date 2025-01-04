package stringutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestRotate 按指定的字符数旋转字符串。
// +1 表示将最右侧的字符移动到最左边
// func Rotate(str string, shift int) string
func TestRotate(t *testing.T) {
	result1 := strutil.Rotate("Hello", 0)
	result2 := strutil.Rotate("Hello", 1)
	result3 := strutil.Rotate("Hello", 2)

	assert.Equal(t, "Hello", result1, "result1")
	assert.Equal(t, "oHell", result2, "result2")
	assert.Equal(t, "loHel", result3, "result3")
}
