package channel

import (
	"context"
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/concurrency"
)

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
