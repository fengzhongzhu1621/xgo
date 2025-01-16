package slice

import (
	"math"
	"sort"
	"testing"

	"github.com/samber/lo"
	"github.com/samber/lo/parallel"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestChunk 按照 size 参数均分 slice
// func Chunk[T any](slice []T, size int) [][]T
func TestChunk(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.Chunk([]int{0, 1, 2, 3, 4, 5}, 2)
		result2 := lo.Chunk([]int{0, 1, 2, 3, 4, 5, 6}, 2)
		result3 := lo.Chunk([]int{}, 2)
		result4 := lo.Chunk([]int{0}, 2)

		is.Equal(result1, [][]int{{0, 1}, {2, 3}, {4, 5}})
		is.Equal(result2, [][]int{{0, 1}, {2, 3}, {4, 5}, {6}})
		is.Equal(result3, [][]int{})
		is.Equal(result4, [][]int{{0}})
		is.PanicsWithValue("Second parameter must be greater than 0", func() {
			lo.Chunk([]int{0}, 0)
		})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Chunk(allStrings, 2)
		is.IsType(nonempty[0], allStrings, "type preserved")

		// appending to a chunk should not affect original array
		originalArray := []int{0, 1, 2, 3, 4, 5}
		result5 := lo.Chunk(originalArray, 2)
		result5[0] = append(result5[0], 6)
		is.Equal(originalArray, []int{0, 1, 2, 3, 4, 5})
	}

	{
		arr := []string{"a", "b", "c", "d", "e"}

		result1 := slice.Chunk(arr, 1)
		result2 := slice.Chunk(arr, 2)
		result3 := slice.Chunk(arr, 3)
		result4 := slice.Chunk(arr, 4)
		result5 := slice.Chunk(arr, 5)
		result6 := slice.Chunk(arr, 6)

		assert.Equal(t, [][]string{{"a"}, {"b"}, {"c"}, {"d"}, {"e"}}, result1)
		assert.Equal(t, [][]string{{"a", "b"}, {"c", "d"}, {"e"}}, result2)
		assert.Equal(t, [][]string{{"a", "b", "c"}, {"d", "e"}}, result3)
		assert.Equal(t, [][]string{{"a", "b", "c", "d"}, {"e"}}, result4)
		assert.Equal(t, [][]string{{"a", "b", "c", "d", "e"}}, result5)
		assert.Equal(t, [][]string{{"a", "b", "c", "d", "e"}}, result6)
	}
}

// TestPartition Partition all slice elements with the evaluation of the given predicate functions.
// func Partition[T any](slice []T, predicates ...func(item T) bool) [][]T
func TestPartition(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		oddEven := func(x int) string {
			if x < 0 {
				return "negative"
			} else if x%2 == 0 {
				return "even"
			}
			return "odd"
		}

		result1 := lo.PartitionBy([]int{-2, -1, 0, 1, 2, 3, 4, 5}, oddEven)
		result2 := lo.PartitionBy([]int{}, oddEven)

		is.Equal(result1, [][]int{{-2, -1}, {0, 2, 4}, {1, 3, 5}})
		is.Equal(result2, [][]int{})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.PartitionBy(allStrings, func(item string) int {
			return len(item)
		})
		is.IsType(nonempty[0], allStrings, "type preserved")
	}

	{
		nums := []int{1, 2, 3, 4, 5}

		result1 := slice.Partition(nums)
		result2 := slice.Partition(nums,
			func(n int) bool { return n%2 == 0 },
		)
		result3 := slice.Partition(nums,
			func(n int) bool {
				return n == 1 || n == 2
			},
			func(n int) bool {
				return n == 2 || n == 3 || n == 4
			},
		)

		assert.Equal(t, [][]int{{1, 2, 3, 4, 5}}, result1)
		assert.Equal(t, [][]int{{2, 4}, {1, 3, 5}}, result2)
		assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, result3)
	}
}

func TestParallelPartitionBy(t *testing.T) {
	is := assert.New(t)

	oddEven := func(x int) string {
		if x < 0 {
			return "negative"
		} else if x%2 == 0 {
			return "even"
		}
		return "odd"
	}

	result1 := parallel.PartitionBy([]int{-2, -1, 0, 1, 2, 3, 4, 5}, oddEven)
	result2 := parallel.PartitionBy([]int{}, oddEven)

	// order
	sort.Slice(result1, func(i, j int) bool {
		return result1[i][0] < result1[j][0]
	})
	for x := range result1 {
		sort.Slice(result1[x], func(i, j int) bool {
			return result1[x][i] < result1[x][j]
		})
	}

	is.ElementsMatch(result1, [][]int{{-2, -1}, {0, 2, 4}, {1, 3, 5}})
	is.Equal(result2, [][]int{})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := parallel.PartitionBy(allStrings, func(item string) int {
		return len(item)
	})
	is.IsType(nonempty[0], allStrings, "type preserved")
}

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

func TestParallelGroupBy(t *testing.T) {
	is := assert.New(t)

	result1 := parallel.GroupBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
		return i % 3
	})

	// order
	for x := range result1 {
		sort.Slice(result1[x], func(i, j int) bool {
			return result1[x][i] < result1[x][j]
		})
	}

	is.EqualValues(len(result1), 3)
	is.EqualValues(result1, map[int][]int{
		0: {0, 3},
		1: {1, 4},
		2: {2, 5},
	})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := parallel.GroupBy(allStrings, func(i string) int {
		return 42
	})
	is.IsType(nonempty[42], allStrings, "type preserved")
}

// TestGroupWith Return a map composed of keys generated from the results of running each element of slice thru iteratee.
// func GroupWith[T any, U comparable](slice []T, iteratee func(T) U) map[U][]T
func TestGroupWith(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.GroupBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
			return i % 3
		})

		is.Equal(len(result1), 3)
		is.Equal(result1, map[int][]int{
			0: {0, 3},
			1: {1, 4},
			2: {2, 5},
		})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.GroupBy(allStrings, func(i string) int {
			return 42
		})
		is.IsType(nonempty[42], allStrings, "type preserved")
	}

	{
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
}

// TestBreak  a slice into two based on a predicate function.
// It starts appending to the second slice after the first element that matches the predicate.
// All elements after the first match are included in the second slice,
// regardless of whether they match the predicate or not.
// 根据谓词函数将一个切片分割为两个。从匹配谓词的第一个元素之后开始将元素追加到第二个切片。
// 第一个匹配元素之后的所有元素（无论是否匹配谓词）都包含在第二个切片中。
// func Break[T any](values []T, predicate func(T) bool) ([]T, []T)
func TestBreak(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	even := func(n int) bool { return n%2 == 0 }

	resultEven, resultAfterFirstEven := slice.Break(nums, even)

	assert.Equal(t, []int{1}, resultEven)
	assert.Equal(t, []int{2, 3, 4, 5}, resultAfterFirstEven)
}
