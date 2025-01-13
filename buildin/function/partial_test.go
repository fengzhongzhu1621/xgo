package function

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/function"
)

// TestCurryFn 编写柯里化函数。这里指的是创建一个柯里化（Currying）函数，柯里化是一种将使用多个参数的函数转换成一系列使用一个参数的函数的技术。
// type CurryFn[T any] func(...T) T
// func (cf CurryFn[T]) New(val T) func(...T) T
func TestCurryFn(t *testing.T) {
	add := func(a, b int) int {
		return a + b
	}

	var addCurry function.CurryFn[int] = func(values ...int) int {
		return add(values[0], values[1])
	}
	// 创建一个偏函数
	add1 := addCurry.New(1)

	result := add1(2)
	fmt.Println(result)

	// Output:
	// 3
}
