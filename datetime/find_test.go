package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
)

// 返回指定日期是该年的第几天
func TestDayOfYear(t *testing.T) {
	date1 := time.Date(2023, 02, 01, 1, 1, 1, 0, time.Local)
	result1 := datetime.DayOfYear(date1)

	date2 := time.Date(2023, 01, 02, 1, 1, 1, 0, time.Local)
	result2 := datetime.DayOfYear(date2)

	date3 := time.Date(2023, 01, 01, 1, 1, 1, 0, time.Local)
	result3 := datetime.DayOfYear(date3)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 31
	// 1
	// 0
}
