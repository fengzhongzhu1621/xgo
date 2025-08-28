package validator

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestChannelIsNotFull(t *testing.T) {
	// 无缓冲通道
	ch1 := make(chan int)
	fmt.Println("ch1 is not full:", ChannelIsNotFull(ch1)) // 输出: true

	// 有缓冲通道
	ch2 := make(chan int, 2)
	fmt.Println("ch2 is not full:", ChannelIsNotFull(ch2)) // 输出: true

	ch2 <- 1
	fmt.Println("ch2 is not full:", ChannelIsNotFull(ch2)) // 输出: true

	ch2 <- 2
	fmt.Println("ch2 is not full:", ChannelIsNotFull(ch2)) // 输出: false

	// 尝试向满的通道发送会导致阻塞
	go func() {
		time.Sleep(time.Second)
		ch2 <- 3
		fmt.Println("Sent 3 to ch2")
	}()

	// 等待发送完成
	time.Sleep(2 * time.Second)
	fmt.Println("Final state of ch2 is not full:", ChannelIsNotFull(ch2)) // 输出: false
}
