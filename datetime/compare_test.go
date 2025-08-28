package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// func Min(t1 time.Time, times ...time.Time) time.Time
func TestMin(t *testing.T) {
	minTime := datetime.Min(
		time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, time.September, 2, 0, 0, 0, 0, time.UTC),
	)

	fmt.Println(minTime)

	// Output:
	// 2024-09-01 00:00:00 +0000 UTC
}

// func Max(t1 time.Time, times ...time.Time) time.Time
func TestMax(t *testing.T) {
	maxTime := datetime.Max(
		time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, time.September, 2, 0, 0, 0, 0, time.UTC),
	)

	fmt.Println(maxTime)

	// Output:
	// 2024-09-02 00:00:00 +0000 UTC
}

// func MaxMin(t1 time.Time, times ...time.Time) (maxTime time.Time, minTime time.Time)
func TestMaxMin(t *testing.T) {
	max, min := datetime.MaxMin(
		time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, time.September, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2024, time.September, 3, 0, 0, 0, 0, time.UTC),
	)

	fmt.Println(max)
	fmt.Println(min)

	// Output:
	// 2024-09-03 00:00:00 +0000 UTC
	// 2024-09-01 00:00:00 +0000 UTC
}

// 获得最小的时间
func TestEarliest(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	a := time.Now()
	b := a.Add(time.Hour)
	result1 := lo.Earliest(a, b)
	result2 := lo.Earliest()

	is.Equal(result1, a)
	is.Equal(result2, time.Time{})
}

func TestEarliestBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	type foo struct {
		bar time.Time
	}

	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	t3 := t1.Add(-time.Hour)
	result1 := lo.EarliestBy([]foo{{t1}, {t2}, {t3}}, func(i foo) time.Time {
		return i.bar
	})
	result2 := lo.EarliestBy([]foo{{t1}}, func(i foo) time.Time {
		return i.bar
	})
	result3 := lo.EarliestBy([]foo{}, func(i foo) time.Time {
		return i.bar
	})

	is.Equal(result1, foo{t3})
	is.Equal(result2, foo{t1})
	is.Equal(result3, foo{})
}

func TestLatest(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	a := time.Now()
	b := a.Add(time.Hour)
	result1 := lo.Latest(a, b)
	result2 := lo.Latest()

	is.Equal(result1, b)
	is.Equal(result2, time.Time{})
}

func TestLatestBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	type foo struct {
		bar time.Time
	}

	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	t3 := t1.Add(-time.Hour)
	result1 := lo.LatestBy([]foo{{t1}, {t2}, {t3}}, func(i foo) time.Time {
		return i.bar
	})
	result2 := lo.LatestBy([]foo{{t1}}, func(i foo) time.Time {
		return i.bar
	})
	result3 := lo.LatestBy([]foo{}, func(i foo) time.Time {
		return i.bar
	})

	is.Equal(result1, foo{t2})
	is.Equal(result2, foo{t1})
	is.Equal(result3, foo{})
}
