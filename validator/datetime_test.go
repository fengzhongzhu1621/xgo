package validator

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
)

// TestIsLeapYear 检查参数`year`是否为闰年。
// func IsLeapYear(year int) bool
func TestIsLeapYear(t *testing.T) {
	result1 := datetime.IsLeapYear(2000)
	result2 := datetime.IsLeapYear(2001)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// true
	// false
}

// Checks if passed time is weekend or not.
// func IsWeekend(t time.Time) bool
func TestIsWeekend(t *testing.T) {
	date1 := time.Date(2023, 06, 03, 0, 0, 0, 0, time.Local)
	date2 := time.Date(2023, 06, 04, 0, 0, 0, 0, time.Local)
	date3 := time.Date(2023, 06, 02, 0, 0, 0, 0, time.Local)

	result1 := datetime.IsWeekend(date1)
	result2 := datetime.IsWeekend(date2)
	result3 := datetime.IsWeekend(date3)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// true
	// true
	// false
}

func TestIsDate(t *testing.T) {
	type args struct {
		sInput string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{"2018-10-10"}, want: true},
		{args: args{"2018/10/10"}, want: false},
		{args: args{`2018\10\10`}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDate(tt.args.sInput); got != tt.want {
				t.Errorf("IsDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTime(t *testing.T) {
	type args struct {
		sInput string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{"2018-10-10 10:56:67"}, want: true},
		{args: args{"105667"}, want: false},
		{args: args{`10-56-67`}, want: false},
		{args: args{"2021-04-07T21:50:50.351153+08:00"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := IsTime(tt.args.sInput); got != tt.want {
				t.Errorf("IsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
