package function

// `AcceptIf`返回一个函数，这个函数与`apply`函数具有相同的签名，并且还包含一个布尔值用以指示成功或失败。
// 一个谓词函数，它接受一个类型为`T`的参数并返回一个布尔值；
// 一个`apply`函数，它也接受一个类型为`T`的参数并返回一个相同类型的修改后的值。
// func AcceptIf[T any](predicate func(T) bool, apply func(T) T) func(T) (T, bool)
import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/function"
)

// 根据条件判断确定是否执行函数
func TestAcceptIf(t *testing.T) {

	adder := function.AcceptIf(
		function.And(
			func(x int) bool {
				return x > 10
			},
			func(x int) bool {
				return x%2 == 0
			}),
		// 条件为 true 时才执行次函数
		func(x int) int {
			return x + 1
		},
	)

	result, ok := adder(20)
	fmt.Println(result) // 21
	fmt.Println(ok)     // true

	result, ok = adder(21)
	fmt.Println(result) // 0
	fmt.Println(ok)     // false
}
