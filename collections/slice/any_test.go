package slice

import (
	"testing"

	"github.com/araujo88/lambda-go/pkg/predicate"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"

	"github.com/stretchr/testify/assert"
)

// TestSome Return true if any of the values in the list pass the predicate function.
// func Some[T any](slice []T, predicate func(index int, item T) bool) bool
func TestSome(t *testing.T) {
	{
		t.Parallel()
		is := assert.New(t)

		result1 := lo.Some([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
		result2 := lo.Some([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
		result3 := lo.Some([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
		result4 := lo.Some([]int{0, 1, 2, 3, 4, 5}, []int{})

		is.True(result1)
		is.True(result2)
		is.False(result3)
		is.False(result4)
	}

	{
		nums := []int{1, 2, 3, 5}
		isEven := func(i, num int) bool {
			return num%2 == 0
		}
		result := slice.Some(nums, isEven)
		assert.Equal(t, true, result)
	}

	{
		tests := []struct {
			name      string
			slice     []int
			predicate func(int) bool
			want      bool
		}{
			{"true for positive match", []int{1, 2, 3}, func(x int) bool { return x == 2 }, true},
			{"false for no match", []int{1, 2, 3}, func(x int) bool { return x == 5 }, false},
			{"empty slice", []int{}, func(x int) bool { return x == 1 }, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := predicate.Any(tt.slice, tt.predicate); got != tt.want {
					t.Errorf("Any() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

func TestSomeBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.SomeBy([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 5
	})

	is.True(result1)

	result2 := lo.SomeBy([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 3
	})

	is.True(result2)

	result3 := lo.SomeBy([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 0
	})

	is.False(result3)

	result4 := lo.SomeBy([]int{}, func(x int) bool {
		return x < 5
	})

	is.False(result4)
}
