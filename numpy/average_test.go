package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

// TestAverage Return average value of numbers. Maybe call RoundToFloat to round result.
// func Average[T constraints.Integer | constraints.Float](numbers ...T) float64
// func RoundToFloat[T constraints.Float | constraints.Integer](x T, n int) float64
func TestAverage(t *testing.T) {
	result1 := mathutil.Average(1, 2)

	avg := mathutil.Average(1.2, 1.25)
	result2 := mathutil.RoundToFloat(avg, 1)
	result3 := mathutil.RoundToFloat(avg, 2)

	fmt.Println(result1) // 1.5
	fmt.Println(result2) // 1.2
	fmt.Println(result3) // 1.23
}
