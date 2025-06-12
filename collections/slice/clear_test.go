package slice

import (
	"fmt"
	"testing"
)

// // The clear built-in function clears maps and slices.
// For maps, clear deletes all entries, resulting in an empty map.
// For slices, clear sets all elements up to the length of the slice
// to the zero value of the respective element type. If the argument
// type is a type parameter, the type parameter's type set must
// contain only map or slice types, and clear performs the operation
// implied by the type argument.
// func clear[T ~[]Type | ~map[Type]Type1](t T)

// 注意：如果 map，slice 为 nil，函数 clear 的执行则是无效操作。
// 如果函数 clear 的入参是 type parameter （类型参数），
// 则类型参数的集合必须仅包含 map 或 slice，函数 clear 则按照类型参数集合中的字段类型，执行相应的操作。

type Data struct {
	User   map[int]string
	Salary map[string]int
}

type Data1 struct {
	User   string
	Salary int
}

func TestClear(t *testing.T) {
	s := []int{1, 2, 3}
	fmt.Printf("len=%d\t s=%+v\n", len(s), s) // len=3      s=[1 2 3]
	clear(s)
	fmt.Printf("len=%d\t s=%+v\n", len(s), s) // len=3      s=[0 0 0]
}

func TestSliceStruct(t *testing.T) {
	d := []Data{
		{
			User:   map[int]string{1: "frank", 2: "lucy"},
			Salary: map[string]int{"frank": 1000, "lucy": 2000},
		},
	}
	fmt.Printf("d=%+v\n", d) // d=[{User:map[1:frank 2:lucy] Salary:map[frank:1000 lucy:2000]}]
	clear(d)
	fmt.Printf("d=%+v\n", d) // d=[{User:map[] Salary:map[]}]

	d1 := []Data1{
		{
			User:   "frank",
			Salary: 1000,
		},
	}
	fmt.Printf("d1=%+v\n", d1) // d1=[{User:frank Salary:1000}]
	clear(d1)
	fmt.Printf("d1=%+v\n", d1) // d1=[{User: Salary:0}]
}
