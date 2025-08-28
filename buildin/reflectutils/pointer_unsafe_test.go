package reflectutils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPointerOf(t *testing.T) {
	x := 42
	y := &x

	// 获取 y 的 reflect.Value
	valueY := reflect.ValueOf(y)

	// 使用 PointerOf 获取 Pointer
	ptr := PointerOf(valueY)

	// 检查指针是否为 nil
	fmt.Println("IsNil:", ptr.IsNil()) // 输出: IsNil: false

	// 获取 uintptr 表示
	fmt.Printf("Uintptr: %x\n", ptr.Uintptr()) // 140001d37d0

	// 修改原始变量并通过指针查看修改
	*(*int)(ptr.p) = 100
	fmt.Println("x after modification:", x) // 输出: x after modification: 100
}
