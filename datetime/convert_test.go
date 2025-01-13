package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
)

func TestTimeStampToLocalString(t *testing.T) {
	timestamp := int64(1633093200)

	str := TimeStampToLocalString(timestamp, time.RFC3339)
	fmt.Println(str)
}

// func FormatTimeToStr(t time.Time, format string, timezone ...string) string
func TestFormatTimeToStr(t *testing.T) {
	t1, _ := time.Parse("2006-01-02 15:04:05", "2021-01-02 16:04:08")

	result1 := datetime.FormatTimeToStr(t1, "yyyy-mm-dd hh:mm:ss")
	result2 := datetime.FormatTimeToStr(t1, "yyyy-mm-dd")
	result3 := datetime.FormatTimeToStr(t1, "dd-mm-yy hh:mm:ss")

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 2021-01-02 16:04:08
	// 2021-01-02
	// 02-01-21 16:04:08
}

// func FormatStrToTime(str, format string, timezone ...string) (time.Time, error)
func TestFormatStrToTime(t *testing.T) {
	result1, _ := datetime.FormatStrToTime("2021-01-02 16:04:08", "yyyy-mm-dd hh:mm:ss")
	result2, _ := datetime.FormatStrToTime("2021-01-02", "yyyy-mm-dd")
	result3, _ := datetime.FormatStrToTime("02-01-21 16:04:08", "dd-mm-yy hh:mm:ss")

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 2021-01-02 16:04:08 +0000 UTC
	// 2021-01-02 00:00:00 +0000 UTC
	// 2021-01-02 16:04:08 +0000 UTC
}

// Return unix timestamp of specified time string, t should be "yyyy-mm-dd hh:mm:ss".
//
//	type theTime struct {
//		unix int64
//	}
//
// func NewFormat(t string) (*theTime, error)
func TestNewFormat(t *testing.T) {
	tm, err := datetime.NewFormat("2022-03-18 17:04:05")
	if err != nil {
		return
	}

	fmt.Println(tm)

	// Output:
	// &{1647594245}
}

// Return unix timestamp of specified iso8601 time string.
//
//	type theTime struct {
//		unix int64
//	}
//
// func NewISO8601(iso8601 string) (*theTime, error)
func TestNewISO8601(t *testing.T) {
	tm, err := datetime.NewISO8601("2006-01-02T15:04:05.999Z")
	if err != nil {
		return
	}
	fmt.Println(tm)

	// Output:
	// &{1136214245}
}

// func (t *theTime) ToUnix() int64
func TestToUnix(t *testing.T) {
	tm := datetime.NewUnixNow()
	fmt.Println(tm.ToUnix())

	// Output:
	// 1647597438
}

// func (t *theTime) ToFormat() string
func TestToFormat(t *testing.T) {
	tm, _ := datetime.NewFormat("2022-03-18 17:04:05")
	fmt.Println(tm.ToFormat())

	// Output:
	// 2022-03-18 17:04:05
}

// func (t *theTime) ToFormatForTpl(tpl string) string
func TestToFormatForTpl(t *testing.T) {
	tm, _ := datetime.NewFormat("2022-03-18 17:04:05")
	ts := tm.ToFormatForTpl("2006/01/02 15:04:05")
	fmt.Println(ts)

	// Output:
	// 2022/03/18 17:04:05
}

// func (t *theTime) ToIso8601() string
func TestToIso8601(t *testing.T) {
	tm, _ := datetime.NewISO8601("2006-01-02T15:04:05.999Z")
	ts := tm.ToIso8601()
	fmt.Println(ts)

	// Output:
	// 2006-01-02T23:04:05+08:00
}
