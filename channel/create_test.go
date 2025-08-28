package channel

import (
	"context"
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/concurrency"
	"github.com/samber/lo"
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

// TestGenerate Creates a channel, then put values into the channel.
// type Channel[T any] struct
// func NewChannel[T any]() *Channel[T]
// func (c *Channel[T]) Generate(ctx context.Context, values ...T) <-chan T
func TestGenerate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := concurrency.NewChannel[int]()

	// 创建一个带有上下文的管道
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
