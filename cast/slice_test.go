package cast

import (
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/samber/lo"

	"github.com/stretchr/testify/assert"
)

func TestChannelToSlice(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 10*time.Millisecond)
	is := assert.New(t)

	ch := lo.SliceToChannel(2, []int{1, 2, 3})
	items := lo.ChannelToSlice(ch)

	is.Equal([]int{1, 2, 3}, items)
}
