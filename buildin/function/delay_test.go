package function

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/function"
)

// TestDebounce 延迟指定间隔后才执行。（非阻塞）
// func Debounce(fn func(), delay time.Duration) (debouncedFn func(), cancelFn func())
// Creates a debounced version of the provided function. The debounced function will only invoke the original function
// after the specified delay has passed since the last time it was invoked.
// It also supports canceling the debounce.
// 创建所提供的函数的去抖动（防抖）版本。该去抖动函数只有在自上次被调用以来经过了指定的延迟时间才会调用原始函数。它还支持取消去抖动操作。
func TestDebounce(t *testing.T) {
	callCount := 0
	fn := func() {
		callCount++
	}

	// 延迟 500 毫秒执行一次
	debouncedFn, _ := function.Debounce(fn, 500*time.Millisecond)

	for i := 0; i < 10; i++ {
		debouncedFn()
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)
	fmt.Println(callCount)
	debouncedFn()

	time.Sleep(1 * time.Second)
	fmt.Println(callCount)

	// Output:
	// 1
	// 2
}

// TestDelay 延迟执行间隔后才执行，使用 sleep 阻塞执行
// Invoke function after delayed time.
// func Delay(delay time.Duration, fn any, args ...any)
func TestDelay(t *testing.T) {
	var print = func(s string) {
		fmt.Println(s)
	}

	function.Delay(2*time.Second, print, "hello")

	// Output:
	// hello
}
