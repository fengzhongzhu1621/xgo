package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/now"
)

func TestConvertTimezone(t *testing.T) {
	t1 := time.Now()
	beijing := time.FixedZone("Beijing Time", 8*3600) // 东八区
	tInBeijing := now.With(t1).In(beijing)            // 转换为北京时间

	fmt.Println(t1)         // 2025-03-13 09:37:05.097814 +0800 CST m=+0.001358293
	fmt.Println(tInBeijing) // 2025-03-13 09:37:05.097814 +0800 Beijing Time
}
