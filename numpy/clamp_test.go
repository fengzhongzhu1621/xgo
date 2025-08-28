package numpy

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// 将一个数值限制在给定的最小值和最大值之间。具体来说，lo.Clamp 函数接受三个参数：
// value：要限制的数值。
// min：允许的最小值。
// max：允许的最大值。
// 函数的返回值将是 value，但如果 value 小于 min，则返回 min；如果 value 大于 max，则返回 max。这样，返回的值将始终在 min 和 max 之间（包括 min 和 max）。
func TestClamp(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Clamp(0, -10, 10)
	result2 := lo.Clamp(-42, -10, 10)
	result3 := lo.Clamp(42, -10, 10)

	is.Equal(result1, 0)
	is.Equal(result2, -10)
	is.Equal(result3, 10)
}
