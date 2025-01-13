package channel

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/concurrency"

	"github.com/stretchr/testify/assert"
)

func TestCreateChannel(t *testing.T) {
	// 创建一个无缓冲的管道
	ch := make(chan int)
	// 创建一个有缓冲的管道
	chBuffered := make(chan int, 10)

	// 发送数据
	ch <- 1
	chBuffered <- 1
	chBuffered <- 2
	chBuffered <- 3

	// 接收数据
	v := <-ch

	assert.Equal(t, v, 1)
}

func TestRange(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 1
	ch <- 1
	for v := range ch {
		assert.Equal(t, v, 1)
	}
}

func TestCollectData(t *testing.T) {
	var m sync.Mutex
	numbers := []int{1, 2, 3, 4, 5}
	sum := 0
	ch := make(chan int)

	// 遍历数组，启动多个 routines
	for _, num := range numbers {
		go func(n int) {
			m.Lock()
			// 求和
			sum += n
			ch <- sum
			m.Unlock()
		}(num)
	}

	// 获得并行计算的结果
	var finalSum int
	for range numbers {
		finalSum = <-ch
		fmt.Println("result is:", finalSum)
	}
	assert.Equal(t, finalSum, 15)
}

func TestIsChannelClosed(t *testing.T) {
	closed := make(chan struct{})
	close(closed)

	withSentValue := make(chan struct{}, 1)
	withSentValue <- struct{}{}

	testCases := []struct {
		Name           string
		Channel        chan struct{}
		ExpectedPanic  bool
		ExpectedClosed bool
	}{
		{
			Name:           "not_closed",
			Channel:        make(chan struct{}),
			ExpectedPanic:  false,
			ExpectedClosed: false,
		},
		{
			Name:           "closed",
			Channel:        closed, // 已关闭的chan
			ExpectedPanic:  false,
			ExpectedClosed: true,
		},
		{
			Name:           "with_sent_value",
			Channel:        withSentValue,
			ExpectedPanic:  true,
			ExpectedClosed: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			testFunc := func() {
				closed := IsChannelClosed(c.Channel)
				assert.EqualValues(t, c.ExpectedClosed, closed)
			}

			if c.ExpectedPanic {
				assert.Panics(t, testFunc)
			} else {
				assert.NotPanics(t, testFunc)
			}
		})
	}
}

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

// TestGenerate Creates a channel, then put values into the channel.
// type Channel[T any] struct
// func NewChannel[T any]() *Channel[T]
// func (c *Channel[T]) Generate(ctx context.Context, values ...T) <-chan T
func TestGenerate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := concurrency.NewChannel[int]()
	intStream := c.Generate(ctx, 1, 2, 3)

	fmt.Println(<-intStream)
	fmt.Println(<-intStream)
	fmt.Println(<-intStream)

	// Output:
	// 1
	// 2
	// 3
}

// Create a channel whose values are taken from another channel with limit number.
// func (c *Channel[T]) Take(ctx context.Context, valueStream <-chan T, number int) <-chan T
func TestTake(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numbers := make(chan int, 5)
	numbers <- 1
	numbers <- 2
	numbers <- 3
	numbers <- 4
	numbers <- 5
	defer close(numbers)

	c := concurrency.NewChannel[int]()
	intStream := c.Take(ctx, numbers, 3)

	for v := range intStream {
		fmt.Println(v)
	}

	// Output:
	// 1
	// 2
	// 3
}

// Create channel, put values into the channel repeatly until cancel the context.
// func (c *Channel[T]) Repeat(ctx context.Context, values ...T) <-chan T
func TestRepeat(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := concurrency.NewChannel[int]()
	// 从重复管道中获取 4 个元素
	intStream := c.Take(ctx, c.Repeat(ctx, 1, 2), 4)

	for v := range intStream {
		fmt.Println(v)
	}

	// Output:
	// 1
	// 2
	// 1
	// 2
}

// Link multiple channels into one channel until cancel the context.
// func (c *Channel[T]) Bridge(ctx context.Context, chanStream <-chan <-chan T) <-chan T
func TestBridge(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := concurrency.NewChannel[int]()
	genVals := func() <-chan <-chan int {
		out := make(chan (<-chan int))
		go func() {
			defer close(out)
			for i := 1; i <= 5; i++ {
				stream := make(chan int, 1)
				stream <- i
				close(stream)
				// 将管道放入管道
				out <- stream
			}
		}()
		// 返回一个管道
		return out
	}

	for v := range c.Bridge(ctx, genVals()) {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

// Merge multiple channels into one channel until cancel the context.
// func (c *Channel[T]) FanIn(ctx context.Context, channels ...<-chan T) <-chan T
func TestFanIn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := concurrency.NewChannel[int]()
	// 创建了一个切片，该切片包含两个元素，每个元素都是一个只读的整数通道（<-chan int）
	channels := make([]<-chan int, 2)

	for i := 0; i < 2; i++ {
		// 从重复管道中获取 2 个元素
		channels[i] = c.Take(ctx, c.Repeat(ctx, i), 2)
	}

	chs := c.FanIn(ctx, channels...)

	for v := range chs {
		fmt.Println(v)
	}

	// 0
	// 1
	// 0
	// 1
}

// Read one or more channels into one channel, will close when any readin channel is closed.
// func (c *Channel[T]) Or(channels ...<-chan T) <-chan T
func TestOr(t *testing.T) {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	c := concurrency.NewChannel[any]()
	<-c.Or(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
	)

	if time.Since(start).Seconds() < 2 {
		fmt.Println("ok")
	}
	// Output:
	// ok
}

// Read a channel into another channel, will close until cancel context.
// func (c *Channel[T]) OrDone(ctx context.Context, channel <-chan T) <-chan T
func TestOrDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := concurrency.NewChannel[int]()
	// 从重复管道中获取 4 个元素
	intStream := c.Take(ctx, c.Repeat(ctx, 1, 2), 4)

	for v := range c.OrDone(ctx, intStream) {
		fmt.Println(v)
		cancel()
	}
	// Output:
	// 1
}

// Split one chanel into two channels, until cancel the context.
// func (c *Channel[T]) Tee(ctx context.Context, in <-chan T) (<-chan T, <-chan T)
func TestTee(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := concurrency.NewChannel[int]()
	intStream := c.Take(ctx, c.Repeat(ctx, 1), 2)

	ch1, ch2 := c.Tee(ctx, intStream)

	for v := range ch1 {
		fmt.Println(v)
		fmt.Println(<-ch2)
	}

	// Output:
	// 1
	// 1
	// 1
	// 1
}
