package cast

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
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

// func MapToSlice[T any, K comparable, V any](aMap map[K]V, iteratee func(K, V) T) []T
func TestMapToSlice(t *testing.T) {
	aMap := map[string]int{"a": 1, "b": 2, "c": 3}
	result := convertor.MapToSlice(aMap, func(key string, value int) string {
		return key + ":" + strconv.Itoa(value)
	})

	fmt.Println(result) //[]string{"a:1", "b:2", "c:3"}
}
