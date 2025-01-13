package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// func Max[T constraints.Integer | constraints.Float](numbers ...T) T
func TestMax(t *testing.T) {
	result1 := mathutil.Max(1, 2, 3)
	result2 := mathutil.Max(1.2, 1.4, 1.1, 1.4)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 3
	// 1.4
}

// func MaxBy[T any](slice []T, comparator func(T, T) bool) T
func TestMaxBy(t *testing.T) {
	result1 := mathutil.MaxBy([]string{"a", "ab", "abc"}, func(v1, v2 string) bool {
		return len(v1) > len(v2)
	})

	result2 := mathutil.MaxBy([]string{"abd", "abc", "ab"}, func(v1, v2 string) bool {
		return len(v1) > len(v2)
	})

	result3 := mathutil.MaxBy([]string{}, func(v1, v2 string) bool {
		return len(v1) > len(v2)
	})

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// abc
	// abd
	//
}

// func Min[T constraints.Integer | constraints.Float](numbers ...T) T
func TestMin(t *testing.T) {
	result1 := mathutil.Min(1, 2, 3)
	result2 := mathutil.Min(1.2, 1.4, 1.1, 1.4)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 1
	// 1.1
}

// func MinBy[T any](slice []T, comparator func(T, T) bool) T
func TestMinBy(t *testing.T) {
	result1 := mathutil.MinBy([]string{"a", "ab", "abc"}, func(v1, v2 string) bool {
		return len(v1) < len(v2)
	})

	result2 := mathutil.MinBy([]string{"ab", "ac", "abc"}, func(v1, v2 string) bool {
		return len(v1) < len(v2)
	})

	result3 := mathutil.MinBy([]string{}, func(v1, v2 string) bool {
		return len(v1) < len(v2)
	})

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// a
	// ab
	//
}
