package channel

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/concurrency"
	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

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
