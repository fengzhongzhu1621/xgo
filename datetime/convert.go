package datetime

import (
	"time"
)

// SecondToTime 把时间戳转换为time.Time格式.
func SecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

// TodayStartTimestamp 返回今天开始时间的 Unix 时间戳（自 1970 年 1月1日 00:00:00 UTC 起经过的秒数）。
func TodayStartTimestamp() int64 {
	now := time.Now()
	year, month, day := now.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	return start.Unix()
}

// TimeStampToLocalString 时间戳转换为本地时间字符串格式
func TimeStampToLocalString(timestamp int64, format string) string {
	if timestamp == 0 {
		return ""
	}
	t := time.Unix(timestamp, 0)
	localTime := t.Local()

	return localTime.Format(format)
}
