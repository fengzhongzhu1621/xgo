package bytesutils

import (
	"errors"
	"io/fs"
	"testing"
	"time"

	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
)

func TestAppendAny(t *testing.T) {
	assert.Eq(t, []byte("123"), byteutil.AppendAny(nil, 123))
	assert.Eq(t, []byte("123"), byteutil.AppendAny([]byte{}, 123))
	assert.Eq(t, []byte("123"), byteutil.AppendAny([]byte("1"), 23))
	assert.Eq(t, []byte("1<nil>"), byteutil.AppendAny([]byte("1"), nil))
	assert.Eq(t, "3600000000000", string(byteutil.AppendAny([]byte{}, timex.OneHour)))

	tests := []struct {
		dst []byte
		v   any
		exp []byte
	}{
		{nil, 123, []byte("123")},
		{[]byte{}, 123, []byte("123")},
		{[]byte("1"), 23, []byte("123")},
		{[]byte("1"), nil, []byte("1<nil>")},
		{[]byte{}, timex.OneHour, []byte("3600000000000")},
		{[]byte{}, int8(123), []byte("123")},
		{[]byte{}, int16(123), []byte("123")},
		{[]byte{}, int32(123), []byte("123")},
		{[]byte{}, int64(123), []byte("123")},
		{[]byte{}, uint(123), []byte("123")},
		{[]byte{}, uint8(123), []byte("123")},
		{[]byte{}, uint16(123), []byte("123")},
		{[]byte{}, uint32(123), []byte("123")},
		{[]byte{}, uint64(123), []byte("123")},
		{[]byte{}, float32(123), []byte("123")},
		{[]byte{}, float64(123), []byte("123")},
		{[]byte{}, "123", []byte("123")},
		{[]byte{}, []byte("123"), []byte("123")},
		{[]byte{}, []int{1, 2, 3}, []byte("[1 2 3]")},
		{[]byte{}, []string{"1", "2", "3"}, []byte("[1 2 3]")},
		{[]byte{}, true, []byte("true")},
		{[]byte{}, fs.ModePerm, []byte("-rwxrwxrwx")},
		{[]byte{}, errors.New("error msg"), []byte("error msg")},
	}

	for _, tt := range tests {
		assert.Eq(t, tt.exp, byteutil.AppendAny(tt.dst, tt.v))
	}

	tim, err := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	assert.NoError(t, err)
	assert.Eq(t, []byte("2019-01-01T00:00:00Z"), byteutil.AppendAny(nil, tim))
}
