package slice

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestToSlicePtr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	str1 := "foo"
	str2 := "bar"
	result1 := lo.ToSlicePtr([]string{str1, str2})

	is.Equal(result1, []*string{&str1, &str2})
}

func TestFromSlicePtr(t *testing.T) {
	is := assert.New(t)

	str1 := "foo"
	str2 := "bar"
	result1 := lo.FromSlicePtr([]*string{&str1, &str2, nil})

	is.Equal(result1, []string{str1, str2, ""})
}

func TestFromSlicePtrOr(t *testing.T) {
	is := assert.New(t)

	str1 := "foo"
	str2 := "bar"
	result1 := lo.FromSlicePtrOr([]*string{&str1, &str2, nil}, "fallback")

	is.Equal(result1, []string{str1, str2, "fallback"})
}
