package cast

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/gookit/goutil/arrutil"
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

func TestToAnySlice(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	in1 := []int{0, 1, 2, 3}
	in2 := []int{}
	out1 := lo.ToAnySlice(in1)
	out2 := lo.ToAnySlice(in2)

	is.Equal([]any{0, 1, 2, 3}, out1)
	is.Equal([]any{}, out2)
}

func TestFromAnySlice(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.NotPanics(func() {
		out1, ok1 := lo.FromAnySlice[string]([]any{"foobar", 42})
		out2, ok2 := lo.FromAnySlice[string]([]any{"foobar", "42"})

		is.Equal([]string{}, out1)
		is.False(ok1)
		is.Equal([]string{"foobar", "42"}, out2)
		is.True(ok2)
	})
}

func TestStringsToInts(t *testing.T) {
	is := assert.New(t)

	ints, err := arrutil.StringsToInts([]string{"1", "2"})
	is.Nil(err)
	is.Equal("[]int{1, 2}", fmt.Sprintf("%#v", ints))

	_, err = arrutil.StringsToInts([]string{"a", "b"})
	is.Error(err)

	ints = arrutil.StringsAsInts([]string{"1", "2"})
	is.Equal("[]int{1, 2}", fmt.Sprintf("%#v", ints))
	is.Nil(arrutil.StringsAsInts([]string{"abc"}))
}

func TestStringsToSlice(t *testing.T) {
	is := assert.New(t)

	as := arrutil.StringsToSlice([]string{"1", "2"})
	is.Equal(`[]interface {}{"1", "2"}`, fmt.Sprintf("%#v", as))
}

func TestAnyToSlice(t *testing.T) {
	is := assert.New(t)

	sl, err := arrutil.AnyToSlice([]int{1, 2})
	is.NoError(err)
	is.Equal("[]interface {}{1, 2}", fmt.Sprintf("%#v", sl))

	_, err = arrutil.AnyToSlice(123)
	is.Error(err)
}
