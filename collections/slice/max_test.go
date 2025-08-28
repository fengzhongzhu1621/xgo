package slice

import (
	"testing"
	"time"

	"github.com/araujo88/lambda-go/pkg/core"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.Max([]int{1, 2, 3})
		result2 := lo.Max([]int{3, 2, 1})
		result3 := lo.Max([]time.Duration{time.Second, time.Minute, time.Hour})
		result4 := lo.Max([]int{})

		is.Equal(result1, 3)
		is.Equal(result2, 3)
		is.Equal(result3, time.Hour)
		is.Equal(result4, 0)
	}

	{
		tests := []struct {
			name string
			a, b int
			want int
		}{
			{"max of 10 and 20", 10, 20, 20},
			{"max of -1 and 1", -1, 1, 1},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := core.Max(tt.a, tt.b); got != tt.want {
					t.Errorf("Max() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

func TestMaxBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.MaxBy([]string{"s1", "string2", "s3"}, func(item string, max string) bool {
		return len(item) > len(max)
	})
	result2 := lo.MaxBy([]string{"string1", "string2", "s3"}, func(item string, max string) bool {
		return len(item) > len(max)
	})
	result3 := lo.MaxBy([]string{}, func(item string, max string) bool {
		return len(item) > len(max)
	})

	is.Equal(result1, "string2")
	is.Equal(result2, "string1")
	is.Equal(result3, "")
}
