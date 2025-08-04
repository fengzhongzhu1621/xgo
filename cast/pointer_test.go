package cast

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestToPtr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.ToPtr([]int{1, 2})

	is.Equal(*result1, []int{1, 2})
}

func TestNil(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	nilFloat64 := lo.Nil[float64]()
	var expNilFloat64 *float64

	nilString := lo.Nil[string]()
	var expNilString *string

	is.Equal(expNilFloat64, nilFloat64)
	is.Nil(nilFloat64)
	is.NotEqual(nil, nilFloat64)

	is.Equal(expNilString, nilString)
	is.Nil(nilString)
	is.NotEqual(nil, nilString)

	is.NotEqual(nilString, nilFloat64)
}

// 零值返回 nil，非零返回参数值的地址
func TestEmptyableToPtr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.Nil(lo.EmptyableToPtr(0))
	is.Nil(lo.EmptyableToPtr(""))
	is.Nil(lo.EmptyableToPtr[[]int](nil))
	is.Nil(lo.EmptyableToPtr[map[int]int](nil))
	is.Nil(lo.EmptyableToPtr[error](nil))

	is.Equal(*lo.EmptyableToPtr(42), 42)
	is.Equal(*lo.EmptyableToPtr("nonempty"), "nonempty")
	is.Equal(*lo.EmptyableToPtr([]int{}), []int{})
	is.Equal(*lo.EmptyableToPtr([]int{1, 2}), []int{1, 2})
	is.Equal(*lo.EmptyableToPtr(map[int]int{}), map[int]int{})
	is.Equal(*lo.EmptyableToPtr(assert.AnError), assert.AnError)
}

func TestFromPtr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	str1 := "foo"
	ptr := &str1

	is.Equal("foo", lo.FromPtr(ptr))
	is.Equal("", lo.FromPtr[string](nil))
	is.Equal(0, lo.FromPtr[int](nil))
	is.Nil(lo.FromPtr[*string](nil))
	is.EqualValues(ptr, lo.FromPtr(&ptr))
}

func TestFromPtrOr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const fallbackStr = "fallback"
	str := "foo"
	ptrStr := &str

	const fallbackInt = -1
	i := 9
	ptrInt := &i

	is.Equal(str, lo.FromPtrOr(ptrStr, fallbackStr))
	is.Equal(fallbackStr, lo.FromPtrOr(nil, fallbackStr))
	is.Equal(i, lo.FromPtrOr(ptrInt, fallbackInt))
	is.Equal(fallbackInt, lo.FromPtrOr(nil, fallbackInt))
}
