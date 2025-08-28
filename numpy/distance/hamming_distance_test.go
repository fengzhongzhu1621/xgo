package distance

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

// TestHammingDistance 汉明距离计算两个字符串之间的汉明距离。汉明距离是对应符号不同的位置数。。
// func HammingDistance(a, b string) (int, error)
func TestHammingDistance(t *testing.T) {
	result1, _ := strutil.HammingDistance("de", "de")
	result2, _ := strutil.HammingDistance("a", "d")

	assert.Equal(t, 0, result1)
	assert.Equal(t, 1, result2)
}
