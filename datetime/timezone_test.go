package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/now"
)

func TestConvertTimezone(t *testing.T) {
	// 使用 In() 方法转换时区
	{
		t1 := time.Now()
		beijing := time.FixedZone("Beijing Time", 8*3600) // 东八区
		tInBeijing := now.With(t1).In(beijing)            // 转换为北京时间

		fmt.Println(t1)         // 2025-03-13 09:37:05.097814 +0800 CST m=+0.001358293
		fmt.Println(tInBeijing) // 2025-03-13 09:37:05.097814 +0800 Beijing Time
	}

	{
		// 加载时区
		loc, _ := time.LoadLocation("America/New_York")

		// 使用特定时区创建时间
		nyTime := time.Date(2023, 11, 10, 12, 0, 0, 0, loc)
		fmt.Println("纽约时间:", nyTime)

		// 转换为UTC时区
		utcTime := nyTime.UTC()
		fmt.Println("UTC时间:", utcTime)

		// 转换为本地时区
		localTime := utcTime.Local()
		fmt.Println("本地时间:", localTime)
	}
}
