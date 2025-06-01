package channel

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/concurrency"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// 将多个通道（<-chan <-chan int）链接成一个通道，直到取消上下文。
// Link multiple channels into one channel until cancel the context.
// func (c *Channel[T]) Bridge(ctx context.Context, chanStream <-chan <-chan T) <-chan T
func TestBridge(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建一个泛型为int的并发通道实例
	c := concurrency.NewChannel[int]()

	// 定义一个匿名函数并赋值给变量 genVals。该函数返回一个类型为 <-chan <-chan int> 的通道，即一个只读的双重只读通道（通道的通道）。
	genVals := func() <-chan <-chan int {
		// 创建外部通道: 创建一个类型为 <-chan int> 的通道，并将其赋值给变量 out。这个通道用于传输内部的只读通道。
		out := make(chan (<-chan int))
		go func() {
			// 确保在 goroutine 结束时关闭外部通道 out，防止资源泄漏
			defer close(out)
			for i := 1; i <= 5; i++ {
				stream := make(chan int, 1)
				stream <- i
				// 关闭 stream 通道，表示不再发送数据
				close(stream)
				// 将 stream 通道发送到外部通道 out 中，使得外部可以通过 out 接收内部的 stream 通道。
				out <- stream
			}
			for i := 6; i <= 10; i++ {
				stream := make(chan int, 1)
				stream <- i
				// 关闭 stream 通道，表示不再发送数据
				close(stream)
				// 将 stream 通道发送到外部通道 out 中，使得外部可以通过 out 接收内部的 stream 通道。
				out <- stream
			}
		}()
		// 返回外部通道 out，供调用者使用。
		return out
	}

	for v := range c.Bridge(ctx, genVals()) {
		fmt.Println(v)
	}
	// Output:
	// 2
	// 1
	// 4
	// 3
	// 5
	// 6
	// 8
	// 7
	// 10
	// 9
}

// 合并多个输入通道（[]<-chan int）的消息到一个缓冲通道中。输出消息没有优先级。当所有的上游通道到达 EOF 时，下游通道关闭。
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
	// 任意一个 chan 关闭后，停止等待
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
		// 上下文取消时，range 不再从 intStream中获取数据了
		cancel()
	}
	// Output:
	// 1
}
