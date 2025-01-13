package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// func Log(n, base float64) float64
func TestLog(t *testing.T) {
	result1 := mathutil.Log(8, 2)
	result2 := mathutil.TruncRound(mathutil.Log(5, 2), 2)
	result3 := mathutil.TruncRound(mathutil.Log(27, 3), 0)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 3
	// 2.32
	// 3
}
