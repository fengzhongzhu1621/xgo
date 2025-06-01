package broadcast

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/concurrency"
)

// Split one chanel into two channels, until cancel the context.
// func (c *Channel[T]) Tee(ctx context.Context, in <-chan T) (<-chan T, <-chan T)
func TestTee(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := concurrency.NewChannel[int]()
	// intStream实际上是一个只发送2个1后关闭的通道。
	intStream := c.Take(ctx, c.Repeat(ctx, 1), 2)

	// Tee会将intStream的数据同时复制到ch1和ch2。ch1和ch2是并行的，Tee会同时向它们发送数据。
	// 由于intStream只发送2个1后关闭，ch1和ch2也会分别收到这2个1，然后关闭。
	ch1, ch2 := c.Tee(ctx, intStream)

	go func() {
		for v := range ch1 {
			fmt.Println("ch1:", v)
		}
	}()
	go func() {
		for v := range ch2 {
			fmt.Println("ch2:", v)
		}
	}()

	time.Sleep(200 * time.Microsecond)
	// Output:
	// ch2: 1
	// ch2: 1
	// ch1: 1
	// ch1: 1
}
