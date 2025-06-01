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
//
// 注意：和 FanOut 的不同之处，Tee 尝试同时向out1和out2发送数据（使用select语句），但是因为 ch1,cha2是阻塞的，
// 如果只消费了 cha1，则 Tee 中 select 的实现会因为 cha2 没有消费而阻塞
// 只有从 ch1 和 cha2 同时消费完了一个数据后，生成者才会写入下一个数据
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
