package channel

import (
	"fmt"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func worker(done chan bool) {
	fmt.Println("Goroutine started")
	// do some work
	fmt.Println("Goroutine finished")
	done <- true
}

// TestWaitDone 测试主进程等待协程执行完成
func TestWaitDone(t *testing.T) {
	done := make(chan bool, 1)
	go worker(done)
	// wait for worker to finish
	<-done
	fmt.Println("Program finished")
}

// TestAsync 在一个 goroutine 中执行一个函数，并将结果返回到一个通道中。
func TestAsync(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 100*time.Millisecond)
	is := assert.New(t)

	sync := make(chan struct{})

	ch := lo.Async(func() int {
		<-sync // 子任务等待
		return 10
	})

	// 启动子任务
	sync <- struct{}{}

	// 等待子任务执行完成
	select {
	case result := <-ch:
		is.Equal(result, 10)
	case <-time.After(time.Millisecond):
		is.Fail("Async should not block")
	}
}

func TestAsyncX(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 100*time.Millisecond)
	is := assert.New(t)

	{
		sync := make(chan struct{})

		ch := lo.Async0(func() {
			<-sync
		})

		sync <- struct{}{}

		select {
		case <-ch:
		case <-time.After(time.Millisecond):
			is.Fail("Async0 should not block")
		}
	}

	{
		sync := make(chan struct{})

		ch := lo.Async1(func() int {
			<-sync
			return 10
		})

		sync <- struct{}{}

		select {
		case result := <-ch:
			is.Equal(result, 10)
		case <-time.After(time.Millisecond):
			is.Fail("Async1 should not block")
		}
	}

	{
		sync := make(chan struct{})

		ch := lo.Async2(func() (int, string) {
			<-sync
			return 10, "Hello"
		})

		sync <- struct{}{}

		select {
		case result := <-ch:
			is.Equal(result, lo.Tuple2[int, string]{A: 10, B: "Hello"})
		case <-time.After(time.Millisecond):
			is.Fail("Async2 should not block")
		}
	}

	{
		sync := make(chan struct{})

		ch := lo.Async3(func() (int, string, bool) {
			<-sync
			return 10, "Hello", true
		})

		sync <- struct{}{}

		select {
		case result := <-ch:
			is.Equal(result, lo.Tuple3[int, string, bool]{A: 10, B: "Hello", C: true})
		case <-time.After(time.Millisecond):
			is.Fail("Async3 should not block")
		}
	}

	{
		sync := make(chan struct{})

		ch := lo.Async4(func() (int, string, bool, float64) {
			<-sync
			return 10, "Hello", true, 3.14
		})

		sync <- struct{}{}

		select {
		case result := <-ch:
			is.Equal(
				result,
				lo.Tuple4[int, string, bool, float64]{A: 10, B: "Hello", C: true, D: 3.14},
			)
		case <-time.After(time.Millisecond):
			is.Fail("Async4 should not block")
		}
	}

	{
		sync := make(chan struct{})

		ch := lo.Async5(func() (int, string, bool, float64, string) {
			<-sync
			return 10, "Hello", true, 3.14, "World"
		})

		sync <- struct{}{}

		select {
		case result := <-ch:
			is.Equal(
				result,
				lo.Tuple5[int, string, bool, float64, string]{
					A: 10,
					B: "Hello",
					C: true,
					D: 3.14,
					E: "World",
				},
			)
		case <-time.After(time.Millisecond):
			is.Fail("Async5 should not block")
		}
	}

	{
		sync := make(chan struct{})

		ch := lo.Async6(func() (int, string, bool, float64, string, int) {
			<-sync
			return 10, "Hello", true, 3.14, "World", 100
		})

		sync <- struct{}{}

		select {
		case result := <-ch:
			is.Equal(
				result,
				lo.Tuple6[int, string, bool, float64, string, int]{
					A: 10,
					B: "Hello",
					C: true,
					D: 3.14,
					E: "World",
					F: 100,
				},
			)
		case <-time.After(time.Millisecond):
			is.Fail("Async6 should not block")
		}
	}
}
