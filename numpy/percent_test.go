package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// TestPercent 计算val占总数的百分比，保留n位小数。
// func Percent(val, total float64, n int) float64
func TestPercent(t *testing.T) {
	result1 := mathutil.Percent(1, 2, 2)
	result2 := mathutil.Percent(0.1, 0.3, 2)
	result3 := mathutil.Percent(-30305, 408420, 2)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 50
	// 33.33
	// -7.42
}
