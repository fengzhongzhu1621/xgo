package function

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/samber/lo"

	"github.com/stretchr/testify/assert"

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

func TestPartial(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	add := func(x float64, y int) string {
		return strconv.Itoa(int(x) + y)
	}
	f := lo.Partial(add, 5)
	is.Equal("15", f(10))
	is.Equal("0", f(-5))
}

func TestPartial1(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	add := func(x float64, y int) string {
		return strconv.Itoa(int(x) + y)
	}
	f := lo.Partial1(add, 5)
	is.Equal("15", f(10))
	is.Equal("0", f(-5))
}

func TestPartial2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	add := func(x float64, y int, z int) string {
		return strconv.Itoa(int(x) + y + z)
	}
	f := lo.Partial2(add, 5)
	is.Equal("24", f(10, 9))
	is.Equal("8", f(-5, 8))
}

func TestPartial3(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	add := func(x float64, y int, z int, a float32) string {
		return strconv.Itoa(int(x) + y + z + int(a))
	}
	f := lo.Partial3(add, 5)
	is.Equal("21", f(10, 9, -3))
	is.Equal("15", f(-5, 8, 7))
}

func TestPartial4(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	add := func(x float64, y int, z int, a float32, b int32) string {
		return strconv.Itoa(int(x) + y + z + int(a) + int(b))
	}
	f := lo.Partial4(add, 5)
	is.Equal("21", f(10, 9, -3, 0))
	is.Equal("14", f(-5, 8, 7, -1))
}

func TestPartial5(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	add := func(x float64, y int, z int, a float32, b int32, c int) string {
		return strconv.Itoa(int(x) + y + z + int(a) + int(b) + c)
	}
	f := lo.Partial5(add, 5)
	is.Equal("26", f(10, 9, -3, 0, 5))
	is.Equal("21", f(-5, 8, 7, -1, 7))
}
