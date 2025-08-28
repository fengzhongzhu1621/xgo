package function

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/function"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
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

	// 仅仅执行一次
	for i := 0; i < 10; i++ {
		debouncedFn()
		time.Sleep(50 * time.Millisecond)
	}

	// 执行第二次
	time.Sleep(1 * time.Second)
	fmt.Println(callCount)
	debouncedFn()

	time.Sleep(1 * time.Second)
	fmt.Println(callCount)

	// Output:
	// 1
	// 2
}

// 为每个不同的键创建一个防抖实例，该实例会延迟调用给定的函数，直到经过指定的等待毫秒数后才执行。
func TestDebounce2(t *testing.T) {
	t.Parallel()

	f1 := func() {
		println("1. Called once after 10ms when func stopped invoking!")
	}
	f2 := func() {
		println("2. Called once after 10ms when func stopped invoking!")
	}
	f3 := func() {
		println("3. Called once after 10ms when func stopped invoking!")
	}

	d1, _ := lo.NewDebounce(10*time.Millisecond, f1)

	// execute 3 times
	//  1. Called once after 10ms when func stopped invoking!
	//	1. Called once after 10ms when func stopped invoking!
	//	1. Called once after 10ms when func stopped invoking!
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			d1()
		}
		time.Sleep(20 * time.Millisecond)
	}

	d2, _ := lo.NewDebounce(10*time.Millisecond, f2)

	// execute once because it is always invoked and only last invoke is worked after 100ms
	// 2. Called once after 10ms when func stopped invoking!
	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			d2()
		}
		time.Sleep(5 * time.Millisecond)
	}

	time.Sleep(10 * time.Millisecond)

	// execute once because it is canceled after 200ms.
	// 3. Called once after 10ms when func stopped invoking!
	d3, cancel := lo.NewDebounce(10*time.Millisecond, f3)
	for i := 0; i < 3; i++ {
		println("i = " + strconv.Itoa(i))
		for j := 0; j < 10; j++ {
			// 当 i == 0 时，不执行 d3
			println("d3")
			d3()
		}
		time.Sleep(20 * time.Millisecond)
		if i == 0 {
			// 取消后续函数的执行
			println("cancel")
			cancel()
		}
	}
}

// 为每个不同的键创建一个防抖实例，该实例会延迟调用给定的函数，直到经过指定的等待毫秒数后才执行。
func TestDebounceBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	mu := sync.Mutex{}
	output := map[int]int{0: 0, 1: 0, 2: 0}

	f1 := func(key int, count int) {
		mu.Lock()
		output[key] += count
		mu.Unlock()
		// fmt.Printf("[key=%d] 1. Called once after 10ms when func stopped invoking!\n", key)
	}
	f2 := func(key int, count int) {
		mu.Lock()
		output[key] += count
		mu.Unlock()
		// fmt.Printf("[key=%d] 2. Called once after 10ms when func stopped invoking!\n", key)
	}
	f3 := func(key int, count int) {
		mu.Lock()
		output[key] += count
		mu.Unlock()
		// fmt.Printf("[key=%d] 3. Called once after 10ms when func stopped invoking!\n", key)
	}

	// execute 3 times
	d1, _ := lo.NewDebounceBy(10*time.Millisecond, f1)
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 3; k++ {
				d1(k)
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	mu.Lock()
	is.EqualValues(output[0], 30)
	is.EqualValues(output[1], 30)
	is.EqualValues(output[2], 30)
	mu.Unlock()

	// execute once because it is always invoked and only last invoke is worked after 100ms
	d2, _ := lo.NewDebounceBy(10*time.Millisecond, f2)
	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 3; k++ {
				d2(k)
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(10 * time.Millisecond)
	mu.Lock()
	is.EqualValues(output[0], 45) // 30 + 5 * 3
	is.EqualValues(output[1], 45)
	is.EqualValues(output[2], 45)
	mu.Unlock()

	// execute once because it is canceled after 200ms.
	d3, cancel := lo.NewDebounceBy(10*time.Millisecond, f3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 3; k++ {
				d3(k)
			}
		}

		time.Sleep(20 * time.Millisecond)
		if i == 0 {
			for k := 0; k < 3; k++ {
				cancel(k)
			}
		}
	}
	mu.Lock()
	is.EqualValues(output[0], 75)
	is.EqualValues(output[1], 75)
	is.EqualValues(output[2], 75)
	mu.Unlock()
}

// TestDelay 延迟执行间隔后才执行，使用 sleep 阻塞执行
// Invoke function after delayed time.
// func Delay(delay time.Duration, fn any, args ...any)
func TestDelay(t *testing.T) {
	print := func(s string) {
		fmt.Println(s)
	}

	function.Delay(2*time.Second, print, "hello")

	// Output:
	// hello
}
