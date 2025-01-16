package slice

import (
	"testing"

	"github.com/samber/lo"

	"github.com/araujo88/lambda-go/pkg/utils"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.Reverse([]int{0, 1, 2, 3, 4, 5})
		result2 := lo.Reverse([]int{0, 1, 2, 3, 4, 5, 6})
		result3 := lo.Reverse([]int{})

		is.Equal(result1, []int{5, 4, 3, 2, 1, 0})
		is.Equal(result2, []int{6, 5, 4, 3, 2, 1, 0})
		is.Equal(result3, []int{})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Reverse(allStrings)
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
		strs := []string{"a", "b", "c", "d"}

		slice.Reverse(strs)
		assert.Equal(t, []string{"d", "c", "b", "a"}, strs)
	}

	{
		slice1 := []int{1, 2, 3}

		reversed := utils.Reverse(slice1)
		assert.Equal(t, []int{3, 2, 1}, reversed)
	}
}
