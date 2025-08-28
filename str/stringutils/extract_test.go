package stringutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestExtractContent 提取源字符串中起始字符串和结束字符串之间的内容。
// func ExtractContent(s, start, end string) []string
func TestExtractContent(t *testing.T) {
	html := `<span>content1</span>aa<span>content2</span>bb<span>content3</span>`

	result := strutil.ExtractContent(html, "<span>", "</span>")

	assert.Equal(t, []string{"content1", "content2", "content3"}, result)
}
