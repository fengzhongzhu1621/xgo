package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
)

// func AddDay(t time.Time, day int64) time.Time
func TestAddDay(t *testing.T) {
	now := time.Now()

	tomorrow := datetime.AddDay(now, 1)
	diff1 := tomorrow.Sub(now)

	yesterday := datetime.AddDay(now, -1)
	diff2 := yesterday.Sub(now)

	fmt.Println(diff1)
	fmt.Println(diff2)

	// Output:
	// 24h0m0s
	// -24h0m0s
}

// func AddHour(t time.Time, hour int64) time.Time
func TestAddHour(t *testing.T) {
	now := time.Now()

	after2Hours := datetime.AddHour(now, 2)
	diff1 := after2Hours.Sub(now)

	before2Hours := datetime.AddHour(now, -2)
	diff2 := before2Hours.Sub(now)

	fmt.Println(diff1)
	fmt.Println(diff2)

	// Output:
	// 2h0m0s
	// -2h0m0s
}

// func AddMinute(t time.Time, minute int64) time.Time
func TestAddMinute(t *testing.T) {
	now := time.Now()

	after2Minutes := datetime.AddMinute(now, 2)
	diff1 := after2Minutes.Sub(now)

	before2Minutes := datetime.AddMinute(now, -2)
	diff2 := before2Minutes.Sub(now)

	fmt.Println(diff1)
	fmt.Println(diff2)

	// Output:
	// 2m0s
	// -2m0s
}

// func AddYear(t time.Time, year int64) time.Time
func TestAddYear(t *testing.T) {
	now := time.Now()

	after1Year := datetime.AddYear(now, 1)
	diff1 := after1Year.Sub(now)

	before1Year := datetime.AddYear(now, -1)
	diff2 := before1Year.Sub(now)

	fmt.Println(diff1)
	fmt.Println(diff2)

	// Output:
	// 8760h0m0s
	// -8760h0m0s
}