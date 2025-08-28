package slice

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/araujo88/lambda-go/pkg/core"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		// min 函数中的参数，如果有浮点型参数，则所有参数都会转换为浮点型参数作比较
		// 浮点型参数-0.0 和 0.0 作为参数，-0.0 小于 0.0；负无穷大，小于任意其它数值；正无穷大，大于任意其它数值。
		c := min(1, 2.0, 3)
		fmt.Printf("%T\t%v\n", c, c) // float64 1
	}

	{
		// min 和 max 的任意参数是 NaN[1]，则返回结果是 NaN ("not-a-number") 值。
		m := min(3.14, math.NaN(), 1.0)
		fmt.Println(m) // NaN
	}

	{
		// 如果 min 函数的入参为字符串类型的参数，则按照字典序返回最小的字符串，如果有空字符串，则返回空字符串。
		// 逐个字节比较，得出最小/最大的字符串，参数可以交换和组合。
		t1 := min("", "foo", "bar")
		fmt.Println(t1) // ""
	}

	{
		result1 := lo.Min([]int{1, 2, 3})
		result2 := lo.Min([]int{3, 2, 1})
		result3 := lo.Min([]time.Duration{time.Second, time.Minute, time.Hour})
		result4 := lo.Min([]int{})

		is.Equal(result1, 1)
		is.Equal(result2, 1)
		is.Equal(result3, time.Second)
		is.Equal(result4, 0)
	}

	{
		tests := []struct {
			name string
			a, b int
			want int
		}{
			{"min of 10 and 20", 10, 20, 10},
			{"min of -1 and 1", -1, 1, -1},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := core.Min(tt.a, tt.b); got != tt.want {
					t.Errorf("Min() = %v, want %v", got, tt.want)
				}
			})
		}
	}
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
