package cast

import (
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"

	"github.com/fengzhongzhu1621/xgo/tests"

	"github.com/stretchr/testify/assert"

	"github.com/duke-git/lancet/v2/convertor"
)

// Convert a collection of elements to a read-only channel.
// func ToChannel[T any](array []T) <-chan T
func TestToChannel(t *testing.T) {
	ch := convertor.ToChannel([]int{1, 2, 3})
	result1 := <-ch
	result2 := <-ch
	result3 := <-ch

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 1
	// 2
	// 3
}

func TestSliceToChannel(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 10*time.Millisecond)
	is := assert.New(t)

	ch := lo.SliceToChannel(2, []int{1, 2, 3})
	items := lo.ChannelToSlice(ch)

	is.Equal([]int{1, 2, 3}, items)
}
