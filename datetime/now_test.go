package datetime

import (
	"fmt"
	"testing"

	"github.com/dromara/carbon/v2"
	"github.com/duke-git/lancet/v2/datetime"
)

// Return unix timestamp of current time
//
//	type theTime struct {
//		unix int64
//	}
//
// func NewUnixNow() *theTime
func TestNewUnixNow(t *testing.T) {
	tm := datetime.NewUnixNow()
	fmt.Println(tm)

	// Output:
	// &{1647597438}
}

func TestNewUnix(t *testing.T) {
	tm := datetime.NewUnix(1647597438)
	fmt.Println(tm)

	// Output:
	// &{1647597438}
}

// func GetNowDate() string
func TestGetNowDate(t *testing.T) {
	currentDate := datetime.GetNowDate()
	fmt.Println(currentDate)

	// Output:
	// 2022-01-28
}

// func GetNowTime() string
func TestGetNowTime(t *testing.T) {
	currentTime := datetime.GetNowTime()
	fmt.Println(currentTime) // 15:57:33

	// Output:
	// 15:57:33
}

// 获取当前时间（系统时区）
func TestGetNowDateTime(t *testing.T) {
	current := datetime.GetNowDateTime()
	// 2022-01-28 15:59:33
	fmt.Println(current)

	now := carbon.Now()
	// 2022-01-28 15:59:33
	fmt.Println("当前时间:", now)
}

// func GetTodayStartTime() string
func TestGetTodayStartTime(t *testing.T) {
	startTime := datetime.GetTodayStartTime()
	fmt.Println(startTime)

	// Output:
	// 2023-06-29 00:00:00
}

// Return the end time of today, format: yyyy-mm-dd 23:59:59.
// func GetTodayEndTime() string
func TestGetTodayEndTime(t *testing.T) {
	endTime := datetime.GetTodayEndTime()
	fmt.Println(endTime)

	// Output:
	// 2023-06-29 23:59:59
}

// Return timestamp of zero hour (timestamp of 00:00).
// func GetZeroHourTimestamp() int64
func TestGetZeroHourTimestamp(t *testing.T) {
	zeroTime := datetime.GetZeroHourTimestamp()
	fmt.Println(zeroTime)

	// Output:
	// 1643299200
}

// Return timestamp of zero hour (timestamp of 23:59).
// func GetNightTimestamp() int64
func TestGetNightTimestamp(t *testing.T) {
	nightTime := datetime.GetNightTimestamp()
	fmt.Println(nightTime)

	// Output:
	// 1643385599
}

// func NowDateOrTime(format string, timezone ...string) string
func TestNowDateOrTime(t *testing.T) {
	result1 := datetime.NowDateOrTime("yyyy-mm-dd hh:mm:ss")
	result2 := datetime.NowDateOrTime("yyyy-mm-dd hh:mm:ss", "EST")

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 2023-07-26 15:01:30
	// 2023-07-26 02:01:30
}

// TestTimestamp Return current second timestamp.
// func Timestamp(timezone ...string) int64
func TestTimestamp(t *testing.T) {
	ts := datetime.Timestamp()

	fmt.Println(ts)

	// Output:
	// 1690363051
}

// func TimestampMilli(timezone ...string) int64
func TestTimestampMilli(t *testing.T) {
	ts := datetime.TimestampMilli()

	fmt.Println(ts)

	// Output:
	// 1690363051331
}

// func TimestampMicro(timezone ...string) int64
func TestTimestampMicro(t *testing.T) {
	ts := datetime.TimestampMicro()

	fmt.Println(ts)

	// Output:
	// 1690363051331784
}

// func TimestampNano(timezone ...string) int64
func TestTimestampNano(t *testing.T) {
	ts := datetime.TimestampNano()

	fmt.Println(ts)

	// Output:
	// 1690363051331788000
}
