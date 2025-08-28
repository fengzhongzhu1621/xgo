package tuple

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/tuple"
)

// TestTuple2 Tuple2 represents a 2 elemnets tuple.
// type Tuple2[A any, B any] struct {
//     FieldA A
//     FieldB B
// }

// func NewTuple2[A any, B any](a A, b B) Tuple2[A, B]
func TestTuple2(t *testing.T) {
	t1 := tuple.NewTuple2(1, 0.1)
	// 1 0.1
	fmt.Printf("%v %v\n", t1.FieldA, t1.FieldB)

	// 1 0.1
	v1, v2 := t1.Unbox()
	fmt.Printf("%v %v\n", v1, v2)
}

// TsetZip2 创建一个元组（Tuple2）切片，其元素对应于给定切片的元素。
// func Zip2[A any, B any](a []A, b []B) []Tuple2[A, B]
// func Unzip2[A any, B any](tuples []Tuple2[A, B]) ([]A, []B)
func TestZip2(t *testing.T) {
	result := tuple.Zip2([]int{1, 2}, []float64{0.1, 0.2})
	// [{1 0.1} {2 0.2}]
	fmt.Println(result)
	v1, v2 := tuple.Unzip2([]tuple.Tuple2[int, float64]{
		{FieldA: 1, FieldB: 0.1},
		{FieldA: 2, FieldB: 0.2},
	})
	// [1 2] [0.1 0.2]
	fmt.Printf("%v %v\n", v1, v2)

	result2 := tuple.Zip2([]int{1, 2}, []float64{0.1})
	// [{1 0.1} {2 0}]
	fmt.Println(result2)
	v3, v4 := tuple.Unzip2(result2)
	// [1 2] [0.1 0]
	fmt.Printf("%v %v\n", v3, v4)
}
