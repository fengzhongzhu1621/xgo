package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/gookit/goutil/fmtutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestDuration(t *testing.T) {
	duration, err := Duration("1h30m")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Duration:", duration) // 输出: Duration: 1h30m0s

	duration, _ = Duration("1s")
	fmt.Println("Duration:", duration) // 输出: Duration: 1s
}

// 获得执行时间
func TestDuration0(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result := lo.Duration(func() { time.Sleep(10 * time.Millisecond) })
	is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
}

// 获得执行时间
func TestDurationX(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result := lo.Duration0(func() { time.Sleep(10 * time.Millisecond) })
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
	}

	{
		a, result := lo.Duration1(func() string { time.Sleep(10 * time.Millisecond); return "a" })
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
		is.Equal("a", a)
	}

	{
		a, b, result := lo.Duration2(
			func() (string, string) { time.Sleep(10 * time.Millisecond); return "a", "b" },
		)
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
		is.Equal("a", a)
		is.Equal("b", b)
	}

	{
		a, b, c, result := lo.Duration3(
			func() (string, string, string) { time.Sleep(10 * time.Millisecond); return "a", "b", "c" },
		)
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
		is.Equal("a", a)
		is.Equal("b", b)
		is.Equal("c", c)
	}

	{
		a, b, c, d, result := lo.Duration4(
			func() (string, string, string, string) { time.Sleep(10 * time.Millisecond); return "a", "b", "c", "d" },
		)
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
		is.Equal("a", a)
		is.Equal("b", b)
		is.Equal("c", c)
		is.Equal("d", d)
	}

	{
		a, b, c, d, e, result := lo.Duration5(func() (string, string, string, string, string) {
			time.Sleep(10 * time.Millisecond)
			return "a", "b", "c", "d", "e"
		})
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
		is.Equal("a", a)
		is.Equal("b", b)
		is.Equal("c", c)
		is.Equal("d", d)
		is.Equal("e", e)
	}

	{
		a, b, c, d, e, f, result := lo.Duration6(
			func() (string, string, string, string, string, string) {
				time.Sleep(10 * time.Millisecond)
				return "a", "b", "c", "d", "e", "f"
			},
		)
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
		is.Equal("a", a)
		is.Equal("b", b)
		is.Equal("c", c)
		is.Equal("d", d)
		is.Equal("e", e)
		is.Equal("f", f)
	}

	{
		a, b, c, d, e, f, g, result := lo.Duration7(
			func() (string, string, string, string, string, string, string) {
				time.Sleep(10 * time.Millisecond)
				return "a", "b", "c", "d", "e", "f", "g"
			},
		)
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
		is.Equal("a", a)
		is.Equal("b", b)
		is.Equal("c", c)
		is.Equal("d", d)
		is.Equal("e", e)
		is.Equal("f", f)
		is.Equal("g", g)
	}

	{
		a, b, c, d, e, f, g, h, result := lo.Duration8(
			func() (string, string, string, string, string, string, string, string) {
				time.Sleep(10 * time.Millisecond)
				return "a", "b", "c", "d", "e", "f", "g", "h"
			},
		)
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
		is.Equal("a", a)
		is.Equal("b", b)
		is.Equal("c", c)
		is.Equal("d", d)
		is.Equal("e", e)
		is.Equal("f", f)
		is.Equal("g", g)
		is.Equal("h", h)
	}

	{
		a, b, c, d, e, f, g, h, i, result := lo.Duration9(
			func() (string, string, string, string, string, string, string, string, string) {
				time.Sleep(10 * time.Millisecond)
				return "a", "b", "c", "d", "e", "f", "g", "h", "i"
			},
		)
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
		is.Equal("a", a)
		is.Equal("b", b)
		is.Equal("c", c)
		is.Equal("d", d)
		is.Equal("e", e)
		is.Equal("f", f)
		is.Equal("g", g)
		is.Equal("h", h)
		is.Equal("i", i)
	}

	{
		a, b, c, d, e, f, g, h, i, j, result := lo.Duration10(
			func() (string, string, string, string, string, string, string, string, string, string) {
				time.Sleep(10 * time.Millisecond)
				return "a", "b", "c", "d", "e", "f", "g", "h", "i", "j"
			},
		)
		is.InEpsilon(10*time.Millisecond, result, float64(2*time.Millisecond))
		is.Equal("a", a)
		is.Equal("b", b)
		is.Equal("c", c)
		is.Equal("d", d)
		is.Equal("e", e)
		is.Equal("f", f)
		is.Equal("g", g)
		is.Equal("h", h)
		is.Equal("i", i)
		is.Equal("j", j)
	}
}

// 获得时间差值
func TestBetween(t *testing.T) {
	{
		// 获得两个时间点的差值（秒数）
		// func BetweenSeconds(t1 time.Time, t2 time.Time) int64
		today := time.Now()
		tomorrow := datetime.AddDay(today, 1)
		yesterday := datetime.AddDay(today, -1)

		result1 := datetime.BetweenSeconds(today, tomorrow)
		result2 := datetime.BetweenSeconds(today, yesterday)

		fmt.Println(result1)
		fmt.Println(result2)

		// Output:
		// 86400
		// -86400
	}

	{
		// 获得两个时间点的差值（天）
		// func DaysBetween(start, end time.Time) int
		start := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, time.September, 10, 0, 0, 0, 0, time.UTC)

		result := datetime.DaysBetween(start, end)

		fmt.Println(result)

		// Output:
		// 9
	}

	{
		time1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		time2 := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
		// 比较两个时间
		if time1.Before(time2) {
			fmt.Println("time1在time2之前")
		}
		if time2.After(time1) {
			fmt.Println("time2在time1之后")
		}
		if time1.Equal(time1) {
			fmt.Println("time1等于time1")
		}
		// 计算时间差
		duration := time2.Sub(time1)
		fmt.Println("时间差:", duration)           // 24h0m0s
		fmt.Println("小时数:", duration.Hours())   // 24
		fmt.Println("分钟数:", duration.Minutes()) // 1440
		fmt.Println("秒数:", duration.Seconds())  // 86400
	}
}

// 生成两个时间点之间的日期时间列表
// func GenerateDatetimesBetween(start, end time.Time, layout string, interval string) ([]string, error)
func TestGenerateDatetimesBetween(t *testing.T) {
	start := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, time.September, 1, 2, 0, 0, 0, time.UTC)

	layout := "2006-01-02 15:04:05"
	interval := "1h"

	result, err := datetime.GenerateDatetimesBetween(start, end, layout, interval)

	fmt.Println(result)
	fmt.Println(err)

	// Output:
	// [2024-09-01 00:00:00 2024-09-01 01:00:00 2024-09-01 02:00:00]
	// <nil>
}

// 返回多长时间之后的文字描述
func TestHowLongAgo(t *testing.T) {
	tests := []struct {
		args int64
		want string
	}{
		{-36, "unknown"},
		{36, "36 secs"},
		{346, "5 mins"},
		{3467, "57 mins"},
		{346778, "4 days"},
		{1200346778, "463 months"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, fmtutil.HowLongAgo(tt.args))
	}
}
