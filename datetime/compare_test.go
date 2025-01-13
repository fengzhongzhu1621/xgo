package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
)

// func Min(t1 time.Time, times ...time.Time) time.Time
func TestMin(t *testing.T) {
	minTime := datetime.Min(time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, time.September, 2, 0, 0, 0, 0, time.UTC))

	fmt.Println(minTime)

	// Output:
	// 2024-09-01 00:00:00 +0000 UTC
}

// func Max(t1 time.Time, times ...time.Time) time.Time
func TestMax(t *testing.T) {
	maxTime := datetime.Max(time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, time.September, 2, 0, 0, 0, 0, time.UTC))

	fmt.Println(maxTime)

	// Output:
	// 2024-09-02 00:00:00 +0000 UTC
}

// func MaxMin(t1 time.Time, times ...time.Time) (maxTime time.Time, minTime time.Time)
func TestMaxMin(t *testing.T) {
	max, min := datetime.MaxMin(time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, time.September, 2, 0, 0, 0, 0, time.UTC), time.Date(2024, time.September, 3, 0, 0, 0, 0, time.UTC))

	fmt.Println(max)
	fmt.Println(min)

	// Output:
	// 2024-09-03 00:00:00 +0000 UTC
	// 2024-09-01 00:00:00 +0000 UTC
}
