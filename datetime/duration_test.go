package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
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

// func BetweenSeconds(t1 time.Time, t2 time.Time) int64
func TestBetweenSeconds(t *testing.T) {
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

// func DaysBetween(start, end time.Time) int
func TestDaysBetween(t *testing.T) {
	start := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, time.September, 10, 0, 0, 0, 0, time.UTC)

	result := datetime.DaysBetween(start, end)

	fmt.Println(result)

	// Output:
	// 9
}

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
