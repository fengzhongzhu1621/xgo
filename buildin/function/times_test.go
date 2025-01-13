package function

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/function"
)

// TestAfter 创建一个函数，当被调用达到n次或n次以上时调用给定的func（函数）。
func TestAfter(t *testing.T) {
	fn := function.After(2, func() {
		fmt.Println("hello")
	})

	// hello
	fn() // 一次不调用
	fn()
	// hello
	fn()
	// hello
	fn()
}

// TestBefore 创建一个函数，当被调用的次数小于n次时调用func（函数）。
// func Before(n int, fn any) func(args ...any) []reflect.Value
func TestBefore(t *testing.T) {
	fn := function.Before(2, func() {
		fmt.Println("hello")
	})

	// hello
	fn()
	// hello
	fn()

	// 调用次数超过 2 次，不再调用
	// not call
	fn()
}
