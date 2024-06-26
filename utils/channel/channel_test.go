package channel

import (
	"fmt"
	"sync"
	"testing"

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
