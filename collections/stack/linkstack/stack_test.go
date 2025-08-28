package linkstack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	st := New[struct{}]()
	st.Push(struct{}{})
	require.Equal(t, 1, st.Size())

	st.Reset()
	require.Equal(t, 0, st.Size())

	v, ok := st.Peek()
	require.False(t, ok)
	require.Equal(t, struct{}{}, v)

	v, ok = st.Pop()
	require.False(t, ok)
	require.Equal(t, struct{}{}, v)

	{
		type foo struct {
			bar string
		}

		st := New[foo]()
		st.Push(foo{bar: "baz"})

		v, ok := st.Peek()
		require.True(t, ok)
		require.Equal(t, foo{bar: "baz"}, v)

		v, ok = st.Pop()
		require.True(t, ok)
		require.Equal(t, foo{bar: "baz"}, v)

		require.Zero(t, st.Size())
	}
}
