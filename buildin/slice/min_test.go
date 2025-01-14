package slice

import (
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Min([]int{1, 2, 3})
	result2 := lo.Min([]int{3, 2, 1})
	result3 := lo.Min([]time.Duration{time.Second, time.Minute, time.Hour})
	result4 := lo.Min([]int{})

	is.Equal(result1, 1)
	is.Equal(result2, 1)
	is.Equal(result3, time.Second)
	is.Equal(result4, 0)
}

func TestMinBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.MinBy([]string{"s1", "string2", "s3"}, func(item string, min string) bool {
		return len(item) < len(min)
	})
	result2 := lo.MinBy([]string{"string1", "string2", "s3"}, func(item string, min string) bool {
		return len(item) < len(min)
	})
	result3 := lo.MinBy([]string{}, func(item string, min string) bool {
		return len(item) < len(min)
	})

	is.Equal(result1, "s1")
	is.Equal(result2, "s3")
	is.Equal(result3, "")
}
