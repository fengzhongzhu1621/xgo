package slice

import (
	"testing"

	"github.com/araujo88/lambda-go/pkg/predicate"
	"github.com/duke-git/lancet/v2/slice"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestEvery Return true if all of the values in the slice pass the predicate function.
// func Every[T any](slice []T, predicate func(index int, item T) bool) bool
func TestEvery(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.Every([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
		result2 := lo.Every([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
		result3 := lo.Every([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
		result4 := lo.Every([]int{0, 1, 2, 3, 4, 5}, []int{})

		is.True(result1)
		is.False(result2)
		is.False(result3)
		is.True(result4)
	}

	{
		nums := []int{1, 2, 3, 5}
		isEven := func(i, num int) bool {
			return num%2 == 0
		}
		result := slice.Every(nums, isEven)
		assert.Equal(t, false, result)
	}

	{
		tests := []struct {
			name      string
			slice     []int
			predicate func(int) bool
			want      bool
		}{
			{"true when all match", []int{2, 4, 6}, func(x int) bool { return x%2 == 0 }, true},
			{"false when one does not match", []int{2, 3, 6}, func(x int) bool { return x%2 == 0 }, false},
			{"empty slice returns true", []int{}, func(x int) bool { return x%2 == 0 }, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := predicate.All(tt.slice, tt.predicate); got != tt.want {
					t.Errorf("All() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

func TestEveryBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.EveryBy([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 5
	})

	is.True(result1)

	result2 := lo.EveryBy([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 3
	})

	is.False(result2)

	result3 := lo.EveryBy([]int{1, 2, 3, 4}, func(x int) bool {
		return x < 0
	})

	is.False(result3)

	result4 := lo.EveryBy([]int{}, func(x int) bool {
		return x < 5
	})

	is.True(result4)
}
