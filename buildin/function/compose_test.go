package function

import (
	"fmt"
	"strings"
	"testing"

	"github.com/duke-git/lancet/v2/function"
)

// TestCompose 从右到左组合函数列表，然后返回组合后的函数。
// func Compose[T any](fnList ...func(...T) T) func(...T) T
func TestCompose(t *testing.T) {
	toUpper := func(strs ...string) string {
		return strings.ToUpper(strs[0])
	}
	toLower := func(strs ...string) string {
		return strings.ToLower(strs[0])
	}
	transform := function.Compose(toUpper, toLower)

	result := transform("aBCde")

	fmt.Println(result)

	// Output:
	// ABCDE
}

// Pipeline takes a list of functions and returns a function whose param will be passed into the functions one by one.
// func Pipeline[T any](funcs ...func(T) T) func(T) T
func TestPipeline(t *testing.T) {
	addOne := func(x int) int {
		return x + 1
	}
	double := func(x int) int {
		return 2 * x
	}
	square := func(x int) int {
		return x * x
	}

	fn := function.Pipeline(addOne, double, square)

	result := fn(2)

	fmt.Println(result)

	// Output:
	// 36
}
