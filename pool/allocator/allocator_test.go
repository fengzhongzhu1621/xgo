package allocator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultAllocator(t *testing.T) {
	bs, free := Malloc(10)
	require.Equal(t, 10, len(bs))
	Free(free)

	// 访问已经释放的字节切片
	value := bs[0]
	require.Equal(t, uint8(0x0), value)

}

func TestClassAllocator(t *testing.T) {
	a := NewClassAllocator()
	bs, free := a.Malloc(10)
	require.Equal(t, 10, len(bs))
	a.Free(free)
}

func TestClassAllocator_InvalidMalloc(t *testing.T) {
	a := NewClassAllocator()
	defer func() {
		require.NotEmpty(t, recover())
	}()
	a.Malloc(-1)
}

func TestClassAllocator_InvalidFree(t *testing.T) {
	a := NewClassAllocator()
	t.Run("free empty slice", func(t *testing.T) {
		defer func() {
			require.NotEmpty(t, recover())
		}()
		// panic
		a.Free(nil)
	})
	t.Run("invalid slice size", func(t *testing.T) {
		defer func() {
			require.NotEmpty(t, recover())
		}()
		a.Free(make([]byte, 9))
	})
}
