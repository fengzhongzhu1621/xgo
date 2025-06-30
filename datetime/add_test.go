package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
)

// func AddDay(t time.Time, day int64) time.Time
func TestAddDay(t *testing.T) {
	now := time.Now()

	{
		tomorrow := datetime.AddDay(now, 1)
		diff1 := tomorrow.Sub(now)

		yesterday := datetime.AddDay(now, -1)
		diff2 := yesterday.Sub(now)

		fmt.Println(diff1)
		fmt.Println(diff2)

		// Output:
		// 24h0m0s
		// -24h0m0s
	}

	{
		// 加1天
		oneDayLater := now.Add(24 * time.Hour)
		fmt.Println("一天后:", oneDayLater)
	}
}

// func AddHour(t time.Time, hour int64) time.Time
func TestAddHour(t *testing.T) {
	now := time.Now()

	{
		after2Hours := datetime.AddHour(now, 2)
		diff1 := after2Hours.Sub(now)

		before2Hours := datetime.AddHour(now, -2)
		diff2 := before2Hours.Sub(now)

		fmt.Println(diff1)
		fmt.Println(diff2)

		// Output:
		// 2h0m0s
		// -2h0m0s
	}

	{
		// 减1小时
		oneHourEarlier := now.Add(-1 * time.Hour)
		fmt.Println("一小时前:", oneHourEarlier)
	}
}

// func AddMinute(t time.Time, minute int64) time.Time
func TestAddMinute(t *testing.T) {
	now := time.Now()

	{
		after2Minutes := datetime.AddMinute(now, 2)
		diff1 := after2Minutes.Sub(now)

		before2Minutes := datetime.AddMinute(now, -2)
		diff2 := before2Minutes.Sub(now)

		fmt.Println(diff1)
		fmt.Println(diff2)

		// Output:
		// 2m0s
		// -2m0s
	}

	{
		// 加30分钟
		thirtyMinutesLater := now.Add(30 * time.Minute)
		fmt.Println("30分钟后:", thirtyMinutesLater)
	}
}

func TestAddMonth(t *testing.T) {
	now := time.Now()

	{
		after1Month := datetime.AddMonth(now, 1)
		diff1 := after1Month.Sub(now)

		before1Month := datetime.AddMonth(now, -1)
		diff2 := before1Month.Sub(now)

		fmt.Println(diff1)
		fmt.Println(diff2)

		// Output:
		// 2678400h0m0s
		// -2678400h0m0s
	}

	{
		oneMonthLater := now.AddDate(0, 1, 0)
		fmt.Println("一个月后:", oneMonthLater)

		threeMonthsAgo := now.AddDate(0, -3, 0)
		fmt.Println("三个月前:", threeMonthsAgo)
	}
}

// func AddYear(t time.Time, year int64) time.Time
func TestAddYear(t *testing.T) {
	now := time.Now()

	{
		after1Year := datetime.AddYear(now, 1)
		diff1 := after1Year.Sub(now)

		before1Year := datetime.AddYear(now, -1)
		diff2 := before1Year.Sub(now)

		fmt.Println(diff1)
		fmt.Println(diff2)

		// Output:
		// 8760h0m0s
		// -8760h0m0s
	}

	{
		oneYearLater := now.AddDate(1, 0, 0)
		fmt.Println("一年后:", oneYearLater)
	}
}

func TestNextMonthFirstDay(t *testing.T) {
	now := time.Now()
	fmt.Println("当前时间:", now)

	{
		// 1. 计算下个月的当前日期
		nextMonthSameDay := now.AddDate(0, 1, 0)

		// 2. 用 time.Date 构造下个月 1 号
		nextMonthFirstDay := time.Date(
			nextMonthSameDay.Year(),
			nextMonthSameDay.Month(),
			1,          // 设置为 1 号
			0, 0, 0, 0, // 时间部分设为 00:00:00.000
			now.Location(), // 保持原时区
		)

		fmt.Println("下个月 1 号:", nextMonthFirstDay)
	}

	{
		// 1. 计算下个月 1 号
		nextMonthFirstDay := now.AddDate(0, 1, -now.Day()+1)

		// 2. 确保时间是 00:00:00.000
		nextMonthFirstDay = time.Date(
			nextMonthFirstDay.Year(),
			nextMonthFirstDay.Month(),
			1,
			0, 0, 0, 0,
			now.Location(),
		)

		fmt.Println("下个月 1 号:", nextMonthFirstDay)
	}
}
