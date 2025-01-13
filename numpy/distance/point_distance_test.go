package distance

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// func PointDistance(x1, y1, x2, y2 float64) float64
func TestPointDistance(t *testing.T) {
	result1 := mathutil.PointDistance(1, 1, 4, 5)

	fmt.Println(result1)

	// Output:
	// 5
}
