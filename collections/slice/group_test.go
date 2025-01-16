package slice

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestCountValuesBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	oddEven := func(v int) bool {
		return v%2 == 0
	}
	length := func(v string) int {
		return len(v)
	}

	result1 := lo.CountValuesBy([]int{}, oddEven)
	result2 := lo.CountValuesBy([]int{1, 2}, oddEven)
	result3 := lo.CountValuesBy([]int{1, 2, 2}, oddEven)
	result4 := lo.CountValuesBy([]string{"foo", "bar", ""}, length)
	result5 := lo.CountValuesBy([]string{"foo", "bar", "bar"}, length)

	is.Equal(map[bool]int{}, result1)
	is.Equal(map[bool]int{true: 1, false: 1}, result2)
	is.Equal(map[bool]int{true: 2, false: 1}, result3)
	is.Equal(map[int]int{0: 1, 3: 2}, result4)
	is.Equal(map[int]int{3: 3}, result5)
}
