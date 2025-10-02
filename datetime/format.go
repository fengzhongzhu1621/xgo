package datetime

import "time"

// Some common used time formats.
const (
	// TimeFormatMinute is accurate to the minute.
	TimeFormatMinute = "%Y%m%d%H%M"
	// TimeFormatHour is accurate to the hour.
	TimeFormatHour = "%Y%m%d%H"
	// TimeFormatDay is accurate to the day.
	TimeFormatDay = "%Y%m%d"
	// TimeFormatMonth is accurate to the month.
	TimeFormatMonth = "%Y%m"
	// TimeFormatYear is accurate to the year.
	TimeFormatYear = "%Y"
)

// TimeUnit is the time unit by which files are split, one of minute/hour/day/month/year.
type TimeUnit string

const (
	// Minute splits by the minute.
	Minute = "minute"
	// Hour splits by the hour.
	Hour = "hour"
	// Day splits by the day.
	Day = "day"
	// Month splits by the month.
	Month = "month"
	// Year splits by the year.
	Year = "year"
)

// Format returns a string preceding with `.`. Use TimeFormatDay as default.
func (t TimeUnit) Format() string {
	var timeFmt string
	switch t {
	case Minute:
		timeFmt = TimeFormatMinute
	case Hour:
		timeFmt = TimeFormatHour
	case Day:
		timeFmt = TimeFormatDay
	case Month:
		timeFmt = TimeFormatMonth
	case Year:
		timeFmt = TimeFormatYear
	default:
		timeFmt = TimeFormatDay
	}
	return "." + timeFmt
}

// RotationGap returns the time.Duration for time unit. Use one day as the default.
func (t TimeUnit) RotationGap() time.Duration {
	switch t {
	case Minute:
		return time.Minute
	case Hour:
		return time.Hour
	case Day:
		return time.Hour * 24
	case Month:
		return time.Hour * 24 * 30
	case Year:
		return time.Hour * 24 * 365
	default:
		return time.Hour * 24
	}
}
