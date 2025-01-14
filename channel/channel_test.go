package channel

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/tests"

	"github.com/samber/lo"

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

	is := assert.New(t)
	generator := func(yield func(int)) {
		yield(0)
		yield(1)
		yield(2)
		yield(3)
	}
	i := 0
	for v := range lo.Generator(2, generator) {
		is.Equal(i, v)
		i++
	}
	is.Equal(i, 4)
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

// Buffer 创建一个包含 n 个元素的切片，这些元素来自通道。返回切片、切片长度、读取时间和通道状态（打开/关闭）。
// 第一个参数是通道，第二个参数是切片的长度。
func TestTake2(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 10*time.Millisecond)
	is := assert.New(t)

	ch := lo.SliceToChannel(2, []int{1, 2, 3})

	items1, length1, _, ok1 := lo.Buffer(ch, 2)
	items2, length2, _, ok2 := lo.Buffer(ch, 2)
	items3, length3, _, ok3 := lo.Buffer(ch, 2)

	is.Equal([]int{1, 2}, items1)
	is.Equal(2, length1)
	is.True(ok1)

	is.Equal([]int{3}, items2)
	is.Equal(1, length2)
	is.False(ok2)

	is.Equal([]int{}, items3)
	is.Equal(0, length3)
	is.False(ok3)
}

// BufferWithTimeout 和Buffer函数类似， 但是增加了一个超时参数， 如果超时，返回已经读取的元素。
func TestTakeWithTimeout(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 200*time.Millisecond)
	is := assert.New(t)

	generator := func(yield func(int)) {
		for i := 0; i < 5; i++ {
			yield(i)
			time.Sleep(10 * time.Millisecond)
		}
	}
	ch := lo.Generator(0, generator)

	items1, length1, _, ok1 := lo.BufferWithTimeout(ch, 20, 15*time.Millisecond)
	is.Equal([]int{0, 1}, items1)
	is.Equal(2, length1)
	is.True(ok1)

	items2, length2, _, ok2 := lo.BufferWithTimeout(ch, 20, 2*time.Millisecond)
	is.Equal([]int{}, items2)
	is.Equal(0, length2)
	is.True(ok2)

	items3, length3, _, ok3 := lo.BufferWithTimeout(ch, 1, 30*time.Millisecond)
	is.Equal([]int{2}, items3)
	is.Equal(1, length3)
	is.True(ok3)

	items4, length4, _, ok4 := lo.BufferWithTimeout(ch, 2, 25*time.Millisecond)
	is.Equal([]int{3, 4}, items4)
	is.Equal(2, length4)
	is.True(ok4)

	items5, length5, _, ok5 := lo.BufferWithTimeout(ch, 3, 25*time.Millisecond)
	is.Equal([]int{}, items5)
	is.Equal(0, length5)
	is.False(ok5)
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

// 合并多个输入通道的消息到一个缓冲通道中。输出消息没有优先级。当所有的上游通道到达 EOF 时，下游通道关闭。
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

	// 创建 3 个只读管道
	is := assert.New(t)
	upstreams := CreateChannels[int](3, 10)
	roupstreams := ChannelsToReadOnly(upstreams)
	for i := range roupstreams {
		go func(i int) {
			upstreams[i] <- 1
			upstreams[i] <- 1
			close(upstreams[i])
		}(i)
	}
	out := lo.FanIn(10, roupstreams...)
	time.Sleep(10 * time.Millisecond)

	// check input channels
	is.Equal(0, len(roupstreams[0]))
	is.Equal(0, len(roupstreams[1]))
	is.Equal(0, len(roupstreams[2]))

	// check channels allocation
	is.Equal(6, len(out))
	is.Equal(10, cap(out))

	// check channels content
	for i := 0; i < 6; i++ {
		msg0, ok0 := <-out
		is.Equal(true, ok0)
		is.Equal(1, msg0)
	}

	// 验证当所有的上游通道到达 EOF 时，下游通道关闭。
	// check it is closed
	time.Sleep(10 * time.Millisecond)
	msg0, ok0 := <-out
	is.Equal(false, ok0)
	is.Equal(0, msg0)
}

// 广播所有上游消息到多个下游通道。当上游通道到达 EOF 时，下游通道关闭。如果任何下游通道已满，广播将暂停。
func TestFanOut(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 100*time.Millisecond)
	is := assert.New(t)

	// 广播所有上游消息到多个下游通道
	upstream := lo.SliceToChannel(10, []int{0, 1, 2, 3, 4, 5})
	rodownstreams := lo.FanOut(3, 10, upstream)
	time.Sleep(10 * time.Millisecond)

	// 验证下游通道的数量
	is.Equal(3, len(rodownstreams))

	// 验证下游通道的容量和长度，并读取下游通道的数据
	for i := range rodownstreams {
		is.Equal(6, len(rodownstreams[i]))
		is.Equal(10, cap(rodownstreams[i]))
		is.Equal([]int{0, 1, 2, 3, 4, 5}, lo.ChannelToSlice(rodownstreams[i]))
	}

	// 验证当上游通道到达 EOF 时，下游通道关闭
	time.Sleep(10 * time.Millisecond)
	for i := range rodownstreams {
		msg, ok := <-rodownstreams[i]
		is.Equal(false, ok)
		is.Equal(0, msg)
	}
}

// Read one or more channels into one channel, will close when any read in channel is closed.
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
