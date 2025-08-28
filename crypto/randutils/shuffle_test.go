package randutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestShuffle 打乱给定字符串的字符顺序。Creates an slice of shuffled values.
// func Shuffle(str string) string
// func Shuffle[T any](slice []T) []T
func TestShuffle(t *testing.T) {
	{
		t.Parallel()
		is := assert.New(t)

		result1 := lo.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		result2 := lo.Shuffle([]int{})

		is.NotEqual(result1, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		is.Equal(result2, []int{})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Shuffle(allStrings)
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
		nums := []int{1, 2, 3, 4, 5}
		result := slice.Shuffle(nums)

		fmt.Println(result)
	}

	{
		result1 := strutil.Shuffle("hello")
		result2 := strutil.Shuffle("hello")

		fmt.Println(result1)
		fmt.Println(result2)
	}
}
