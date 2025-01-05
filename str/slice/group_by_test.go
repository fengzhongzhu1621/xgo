package slice

import (
	"math"
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestGroupBy Iterates over elements of the slice, each element will be group by criteria, returns two slices.
// func GroupBy[T any](slice []T, groupFn func(index int, item T) bool) ([]T, []T)
func TestGroupBy(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	isEven := func(i, num int) bool {
		return num%2 == 0
	}

	even, odd := slice.GroupBy(nums, isEven)

	assert.Equal(t, []int{2, 4}, even)
	assert.Equal(t, []int{1, 3, 5}, odd)
}

// TestGroupWith Return a map composed of keys generated from the results of running each element of slice thru iteratee.
// func GroupWith[T any, U comparable](slice []T, iteratee func(T) U) map[U][]T
func TestGroupWith(t *testing.T) {
	nums := []float64{6.1, 4.2, 6.3}

	floor := func(num float64) float64 {
		return math.Floor(num)
	}

	result := slice.GroupWith(nums, floor) //map[float64][]float64

	assert.Equal(t, map[float64][]float64{
		4.0: {4.2},
		6.0: {6.1, 6.3},
	}, result)
}
