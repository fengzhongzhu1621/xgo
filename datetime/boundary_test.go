package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/now"

	"github.com/duke-git/lancet/v2/datetime"
)

// 当天零点
func TestBeginningOfDay(t *testing.T) {
	t1 := time.Now()
	startOfDay := now.With(t1).BeginningOfDay()
	fmt.Println(startOfDay) // 2025-03-13 00:00:00 +0800 CST
}

// func BeginOfMinute(t time.Time) time.Time
func TestBeginOfMinute(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.BeginOfMinute(input)

	fmt.Println(result) // 2023-01-08 18:50:00 +0000 UTC
}

// func BeginOfHour(t time.Time) time.Time
func TestBeginOfHour(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.BeginOfHour(input)

	fmt.Println(result) // 2023-01-08 18:00:00 +0000 UTC
}

// func BeginOfDay(t time.Time) time.Time
func TestBeginOfDay(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.BeginOfDay(input)

	fmt.Println(result) // 2023-01-08 00:00:00 +0000 UTC
}

// 返回一周的开始时间（一周从星期日开始）。
// func BeginOfWeek(t time.Time, beginFrom ...time.Weekday) time.Time
func TestBeginOfWeek(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.BeginOfWeek(input)

	fmt.Println(result) // 2023-01-08 00:00:00 +0000 UTC
}

// func BeginOfMonth(t time.Time) time.Time
func TestBeginOfMonth(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.BeginOfMonth(input)

	fmt.Println(result) // 2023-01-01 00:00:00 +0000 UTC
}

// func BeginOfYear(t time.Time) time.Time
func TestBeginOfYear(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.BeginOfYear(input)

	fmt.Println(result) // 2023-01-01 00:00:00 +0000 UTC
}

// func EndOfMinute(t time.Time) time.Time
func TestEndOfMinute(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.EndOfMinute(input)

	fmt.Println(result) // 2023-01-08 18:50:59.999999999 +0000 UTC
}

// func EndOfHour(t time.Time) time.Time
func TestEndOfHour(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.EndOfHour(input)

	fmt.Println(result) // 2023-01-08 18:59:59.999999999 +0000 UTC
}

// func EndOfDay(t time.Time) time.Time
func TestEndOfDay(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.EndOfDay(input)

	fmt.Println(result) // 2023-01-08 23:59:59.999999999 +0000 UTC
}

// func EndOfWeek(t time.Time, endWith ...time.Weekday) time.Time
func TestEndOfWeek(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.EndOfWeek(input)

	fmt.Println(result) // 2023-01-14 23:59:59.999999999 +0000 UTC
}

// func EndOfMonth(t time.Time) time.Time
func TestEndOfMonth(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.EndOfMonth(input)

	fmt.Println(result) // 2023-01-31 23:59:59.999999999 +0000 UTC
}

// 当月最后一天的23:59:59.999999999
func TestEndOfMonth2(t *testing.T) {
	t1 := time.Now()
	endOfMonth := now.With(t1).EndOfMonth()
	fmt.Println(endOfMonth) // 2025-03-31 23:59:59.999999999 +0800 CST
}

// func EndOfYear(t time.Time) time.Time
func TestEndOfYear(t *testing.T) {
	input := time.Date(2023, 1, 8, 18, 50, 10, 100, time.UTC)
	result := datetime.EndOfYear(input)

	fmt.Println(result)

	// Output:
	// 2023-12-31 23:59:59.999999999 +0000 UTC
}
