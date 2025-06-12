package slice

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/araujo88/lambda-go/pkg/tuple"
)

// Helper function min to calculate the minimum of two lengths
func min2[T int](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// tuple 包提供了一种通用的 tuple 数据结构，可用于处理成对数据。
// 当你需要从函数中返回多个值或将相关数据分组时，这尤其方便。
// 使用 Zip 函数将两个片段合并为一个元组片段的示例
func TestZip(t *testing.T) {
	names := []string{"Alice", "Bob", "Charlie"}
	ages := []int{25, 30, 35}

	pairs := tuple.Zip(names, ages)
	for _, pair := range pairs {
		fmt.Printf("%s is %d years old\n", pair.First, pair.Second)
	}

	tests := []struct {
		name   string
		slice1 []int
		slice2 []string
		want   []tuple.Tuple[int, string]
	}{
		{
			name:   "equal length slices",
			slice1: []int{1, 2, 3},
			slice2: []string{"one", "two", "three"},
			want: []tuple.Tuple[int, string]{
				{First: 1, Second: "one"},
				{First: 2, Second: "two"},
				{First: 3, Second: "three"},
			},
		},
		{
			name:   "first slice longer",
			slice1: []int{1, 2, 3, 4},
			slice2: []string{"one", "two", "three"},
			want: []tuple.Tuple[int, string]{
				{First: 1, Second: "one"},
				{First: 2, Second: "two"},
				{First: 3, Second: "three"},
			},
		},
		{
			name:   "second slice longer",
			slice1: []int{1, 2},
			slice2: []string{"one", "two", "three", "four"},
			want: []tuple.Tuple[int, string]{
				{First: 1, Second: "one"},
				{First: 2, Second: "two"},
			},
		},
		{
			name:   "empty first slice",
			slice1: []int{},
			slice2: []string{"one", "two", "three"},
			want:   []tuple.Tuple[int, string]{},
		},
		{
			name:   "empty second slice",
			slice1: []int{1, 2, 3},
			slice2: []string{},
			want:   []tuple.Tuple[int, string]{},
		},
		{
			name:   "both slices empty",
			slice1: []int{},
			slice2: []string{},
			want:   []tuple.Tuple[int, string]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tuple.Zip(tt.slice1, tt.slice2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Zip() = %v, want %v", got, tt.want)
			}
		})
	}
}
