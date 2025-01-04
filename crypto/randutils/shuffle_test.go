package randutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
)

// TestShuffle 打乱给定字符串的字符顺序。
// func Shuffle(str string) string
func TestShuffle(t *testing.T) {
	result1 := strutil.Shuffle("hello")
	result2 := strutil.Shuffle("hello")

	fmt.Println(result1)
	fmt.Println(result2)
}
