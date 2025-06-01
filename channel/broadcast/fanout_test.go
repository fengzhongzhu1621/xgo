package broadcast

import (
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

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
