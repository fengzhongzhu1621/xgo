package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/duke-git/lancet/v2/datetime"
	"github.com/jinzhu/now"
)

func TestTimestampToTime(t *testing.T) {
	// 获取时间戳
	now1 := time.Now()
	unix := now1.Unix() // 秒级时间戳

	// 从时间戳恢复时间
	t1 := time.Unix(unix, 0)
	fmt.Println("从时间戳恢复:", t1)
}

// 时间戳转换为本地时间字符串格式
func TestTimeStampToLocalString(t *testing.T) {
	timestamp := int64(1633093200)

	str := TimeStampToLocalString(timestamp, time.RFC3339)
	fmt.Println(str)
}

// time.Time 转换为字符串格式
// func FormatTimeToStr(t time.Time, format string, timezone ...string) string
func TestFormatTimeToStr(t *testing.T) {
	{
		// 字符串转换为 time.Time 格式
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

	{
		// func (t *theTime) ToFormat() string
		tm, _ := datetime.NewFormat("2022-03-18 17:04:05")
		fmt.Println(tm.ToFormat())

		// Output:
		// 2022-03-18 17:04:05
	}

	{
		// func (t *theTime) ToFormatForTpl(tpl string) string
		tm, _ := datetime.NewFormat("2022-03-18 17:04:05")
		ts := tm.ToFormatForTpl("2006/01/02 15:04:05")
		fmt.Println(ts)

		// Output:
		// 2022/03/18 17:04:05
	}

	{
		now1 := carbon.Now()
		formatted := now1.Format("Y-m-d H:i:s") // 系统时区
		fmt.Println("格式化后的时间:", formatted)
	}

	{
		// 原生实现
		now1 := time.Now()
		// 常用格式（当前时区）
		fmt.Println(now1.Format("2006-01-02"))            // 2025-04-25
		fmt.Println(now1.Format("20060102"))              // 2025-04-25
		fmt.Println(now1.Format("2006-01-02 15:04:05"))   // 2025-04-25 09:53:38
		fmt.Println(now1.Format("15:04:05"))              // 09:53:38
		fmt.Println(now1.Format("2006年01月02日 15时04分05秒")) // 2025年04月25日 09时53分38秒

		// 预定义格式（当前时区）
		fmt.Println("RFC3339: " + now1.Format(time.RFC3339)) // 2025-04-25T09:53:38+08:00
		fmt.Println("RFC1123: " + now1.Format(time.RFC1123)) // Fri, 25 Apr 2025 09:53:38 CST
	}
}

// 字符串转换为 time.Time 格式
// func FormatStrToTime(str, format string, timezone ...string) (time.Time, error)
func TestFormatStrToTime(t *testing.T) {
	result1, _ := datetime.FormatStrToTime("2021-01-02 16:04:08", "yyyy-mm-dd hh:mm:ss")
	result2, _ := datetime.FormatStrToTime("2021-01-02", "yyyy-mm-dd")
	result3, _ := datetime.FormatStrToTime("02-01-21 16:04:08", "dd-mm-yy hh:mm:ss")

	fmt.Println(result1) // 2021-01-02 16:04:08 +0000 UTC
	fmt.Println(result2) // 2021-01-02 00:00:00 +0000 UTC
	fmt.Println(result3) // 2021-01-02 16:04:08 +0000 UTC

	{
		// 不推荐使用
		t1, _ := now.Parse("2023-10-05")       // 自动识别日期格式
		t2, _ := now.Parse("2023/10/05 14:30") // 支持带时间的字符串
		t3, _ := now.Parse("2023-10-05 14:30")
		t4, _ := now.Parse("2023-10-05T01:02:03Z04:00")
		fmt.Println(t1) // 2023-10-05 00:00:00 +0800 CST
		fmt.Println(t2) // 0001-01-01 00:00:00 +0000 UTC
		fmt.Println(t3) // 2023-10-05 14:30:00 +0800 CST
		fmt.Println(t4) // 0001-01-01 00:00:00 +0000 UTC
	}

	{
		// 解析自定义格式（默认0时区）
		t1, _ := time.Parse("2006-01-02", "2023-11-10")
		fmt.Println(t1) // 2023-11-10 00:00:00 +0000 UTC
		// 解析RFC3339格式
		t2, _ := time.Parse(time.RFC3339, "2023-11-10T14:30:45+08:00")
		fmt.Println(t2) // 2023-11-10 14:30:45 +0800 CST
	}
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

	fmt.Println(tm) // &{1647594245}
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

// func (t *theTime) ToIso8601() string
func TestToIso8601(t *testing.T) {
	tm, _ := datetime.NewISO8601("2006-01-02T15:04:05.999Z")
	ts := tm.ToIso8601()
	fmt.Println(ts)

	// Output:
	// 2006-01-02T23:04:05+08:00
}
