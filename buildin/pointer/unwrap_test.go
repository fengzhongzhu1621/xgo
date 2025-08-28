package pointer

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/pointer"
)

// TestUnwrap Returns the value from the pointer.
// func Unwrap[T any](p *T) T
func TestUnwrap(t *testing.T) {
	a := 123
	b := "abc"

	result1 := pointer.Unwrap(&a)
	result2 := pointer.Unwrap(&b)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 123
	// abc
}

// TestUnwrapOr Returns the value from the pointer or fallback if the pointer is nil.
// UnwrapOr[T any](p *T, fallback T) T
func TestUnwrapOr(t *testing.T) {
	a := 123
	b := "abc"

	var c *int
	var d *string

	result1 := pointer.UnwrapOr(&a, 456)
	result2 := pointer.UnwrapOr(&b, "abc")
	result3 := pointer.UnwrapOr(c, 456)   // 取默认值
	result4 := pointer.UnwrapOr(d, "def") // 取默认值

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)

	// Output:
	// 123
	// abc
	// 456
	// def
}
