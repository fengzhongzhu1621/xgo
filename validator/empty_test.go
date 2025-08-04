package validator

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	type test struct{}

	is.Empty(lo.Empty[string]())
	is.Empty(lo.Empty[int64]())
	is.Empty(lo.Empty[test]())
	is.Empty(lo.Empty[chan string]())
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	type test struct {
		foobar string
	}

	is.True(lo.IsEmpty(""))
	is.False(lo.IsEmpty("foo"))
	is.True(lo.IsEmpty[int64](0))
	is.False(lo.IsEmpty[int64](42))
	is.True(lo.IsEmpty(test{foobar: ""}))
	is.False(lo.IsEmpty(test{foobar: "foo"}))
}

func TestIsNotEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	type test struct {
		foobar string
	}

	is.False(lo.IsNotEmpty(""))
	is.True(lo.IsNotEmpty("foo"))
	is.False(lo.IsNotEmpty[int64](0))
	is.True(lo.IsNotEmpty[int64](42))
	is.False(lo.IsNotEmpty(test{foobar: ""}))
	is.True(lo.IsNotEmpty(test{foobar: "foo"}))
}
