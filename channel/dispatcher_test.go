package channel

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// ChannelDispatcher 将信息从输入通道消息分发到 N 个子通道中。当输入通道关闭时，这个关闭事件会被传播到所有的子通道，也就是说，所有的子通道也会被关闭。
// 这些通道可以有一个固定的缓冲容量，或者当 cap（容量）为 0 时，它们是无缓冲的。
// 作用：启用多个消费者通道从输入通道消费数据并处理，注意并非广播模式

// 分配策略
// lo.DispatchingStrategyRoundRobin: 使用轮询策略将消息分发到子通道中。
// lo.DispatchingStrategyRandom: 使用随机策略将消息分发到子通道中。
// lo.DispatchingStrategyWeightedRandom: 使用加权随机策略将消息分发到子通道中。
// lo.DispatchingStrategyFirst: 分发消息到第一个非满的子通道中。
// lo.DispatchingStrategyLeast: 分发消息到最空的子通道中。
// lo.DispatchingStrategyMost: 分发消息到最满的子通道中。
func TestChannelDispatcher(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 100*time.Millisecond)
	is := assert.New(t)

	ch := make(chan int, 10)

	ch <- 0
	ch <- 1
	ch <- 2
	ch <- 3

	is.Equal(4, len(ch))

	// 创建 5 个子通道，缓存是 10
	children := lo.ChannelDispatcher(ch, 5, 10, lo.DispatchingStrategyRoundRobin[int])
	time.Sleep(10 * time.Millisecond)

	// check channels allocation
	is.Equal(5, len(children))
	is.Equal(10, cap(children[0]))
	is.Equal(10, cap(children[1]))
	is.Equal(10, cap(children[2]))
	is.Equal(10, cap(children[3]))
	is.Equal(10, cap(children[4]))

	is.Equal(1, len(children[0]))
	is.Equal(1, len(children[1]))
	is.Equal(1, len(children[2]))
	is.Equal(1, len(children[3]))
	is.Equal(0, len(children[4])) // 空闲一个通道

	// check channels content
	is.Equal(0, len(ch))

	msg0, ok0 := <-children[0]
	is.Equal(ok0, true)
	is.Equal(msg0, 0)

	msg1, ok1 := <-children[1]
	is.Equal(ok1, true)
	is.Equal(msg1, 1)

	msg2, ok2 := <-children[2]
	is.Equal(ok2, true)
	is.Equal(msg2, 2)

	msg3, ok3 := <-children[3]
	is.Equal(ok3, true)
	is.Equal(msg3, 3)

	// msg4, ok4 := <-children[4]
	// is.Equal(ok4, false)
	// is.Equal(msg4, 0)
	// is.Nil(children[4])

	// 关闭输入通道
	// check it is closed
	close(ch)
	time.Sleep(10 * time.Millisecond)
	is.Panics(func() {
		ch <- 42
	})

	msg0, ok0 = <-children[0]
	is.Equal(ok0, false)
	is.Equal(msg0, 0)

	msg1, ok1 = <-children[1]
	is.Equal(ok1, false)
	is.Equal(msg1, 0)

	msg2, ok2 = <-children[2]
	is.Equal(ok2, false)
	is.Equal(msg2, 0)

	msg3, ok3 = <-children[3]
	is.Equal(ok3, false)
	is.Equal(msg3, 0)

	msg4, ok4 := <-children[4]
	is.Equal(ok4, false)
	is.Equal(msg4, 0)

	// unbuffered channels
	children = lo.ChannelDispatcher(ch, 5, 0, lo.DispatchingStrategyRoundRobin[int])
	is.Equal(0, cap(children[0]))
}

// 测试分配策略
func TestDispatchingStrategyRoundRobin(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 10*time.Millisecond)
	is := assert.New(t)

	children := CreateChannels[int](3, 2)
	rochildren := ChannelsToReadOnly(children)
	defer CloseChannels(children)

	is.Equal(0, lo.DispatchingStrategyRoundRobin(42, 0, rochildren))
	is.Equal(1, lo.DispatchingStrategyRoundRobin(42, 1, rochildren))
	is.Equal(2, lo.DispatchingStrategyRoundRobin(42, 2, rochildren))
	is.Equal(0, lo.DispatchingStrategyRoundRobin(42, 3, rochildren))
}

func TestDispatchingStrategyRandom(t *testing.T) {
	tests.TestWithTimeout(t, 10*time.Millisecond)
	is := assert.New(t)

	// with this seed, the order of random channels are: 1 - 0
	rand.Seed(14)

	children := CreateChannels[int](2, 2)
	rochildren := ChannelsToReadOnly(children)
	defer CloseChannels(children)

	for i := 0; i < 2; i++ {
		children[1] <- i
	}

	is.Equal(0, lo.DispatchingStrategyRandom(42, 0, rochildren))
}

func TestDispatchingStrategyWeightedRandom(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 10*time.Millisecond)
	is := assert.New(t)

	children := CreateChannels[int](2, 2)
	rochildren := ChannelsToReadOnly(children)
	defer CloseChannels(children)

	dispatcher := lo.DispatchingStrategyWeightedRandom[int]([]int{0, 42})

	is.Equal(1, dispatcher(42, 0, rochildren))
	children[0] <- 0
	is.Equal(1, dispatcher(42, 0, rochildren))
	children[1] <- 1
	is.Equal(1, dispatcher(42, 0, rochildren))
}

func TestDispatchingStrategyFirst(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 10*time.Millisecond)
	is := assert.New(t)

	children := CreateChannels[int](2, 2)
	rochildren := ChannelsToReadOnly(children)
	defer CloseChannels(children)

	is.Equal(0, lo.DispatchingStrategyFirst(42, 0, rochildren))
	children[0] <- 0
	is.Equal(0, lo.DispatchingStrategyFirst(42, 0, rochildren))
	children[0] <- 1
	is.Equal(1, lo.DispatchingStrategyFirst(42, 0, rochildren))
}

func TestDispatchingStrategyLeast(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 10*time.Millisecond)
	is := assert.New(t)

	children := CreateChannels[int](2, 2)
	rochildren := ChannelsToReadOnly(children)
	defer CloseChannels(children)

	is.Equal(0, lo.DispatchingStrategyLeast(42, 0, rochildren))
	children[0] <- 0
	is.Equal(1, lo.DispatchingStrategyLeast(42, 0, rochildren))
	children[1] <- 0
	is.Equal(0, lo.DispatchingStrategyLeast(42, 0, rochildren))
	children[0] <- 1
	is.Equal(1, lo.DispatchingStrategyLeast(42, 0, rochildren))
	children[1] <- 1
	is.Equal(0, lo.DispatchingStrategyLeast(42, 0, rochildren))
}

func TestDispatchingStrategyMost(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 10*time.Millisecond)
	is := assert.New(t)

	children := CreateChannels[int](2, 2)
	rochildren := ChannelsToReadOnly(children)
	defer CloseChannels(children)

	is.Equal(0, lo.DispatchingStrategyMost(42, 0, rochildren))
	children[0] <- 0
	is.Equal(0, lo.DispatchingStrategyMost(42, 0, rochildren))
	children[1] <- 0
	is.Equal(0, lo.DispatchingStrategyMost(42, 0, rochildren))
	children[0] <- 1
	is.Equal(0, lo.DispatchingStrategyMost(42, 0, rochildren))
	children[1] <- 1
	is.Equal(0, lo.DispatchingStrategyMost(42, 0, rochildren))
}
