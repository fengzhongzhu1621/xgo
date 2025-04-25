package datetime

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeDate(t *testing.T) {
	// 使用日期创建
	specificTime := time.Date(2023, time.November, 10, 23, 0, 0, 0, time.UTC)

	fmt.Println("特定时间:", specificTime)
}
