package validator

import (
	"fmt"
	"regexp"
	"time"
)

type DateTimeFieldType string

const (
	// TimeWithoutLocationType the common date time type which is used by front end and api
	TimeWithoutLocationType DateTimeFieldType = "time_without_location"
	// TimeWithLocationType the date time type compatible for values from db which is marshaled with time zone
	TimeWithLocationType DateTimeFieldType = "time_with_location"
	InvalidDateTimeType  DateTimeFieldType = "invalid"
)

const (
	// mailPattern     = `^[a-z0-9A-Z]+([\-_\.][a-z0-9A-Z]+)*@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)*\.)+[a-zA-Z]{2,4}$`
	datePattern             = `^[0-9]{4}[\-]{1}[0-9]{2}[\-]{1}[0-9]{2}$`
	dateTimePattern         = `^[0-9]{4}[\-]{1}[0-9]{2}[\-]{1}[0-9]{2}[\s]{1}[0-9]{2}[\:]{1}[0-9]{2}[\:]{1}[0-9]{2}$`
	timeWithLocationPattern = `^[0-9]{4}[\-]{1}[0-9]{2}[\-]{1}[0-9]{2}[T]{1}[0-9]{2}[\:]{1}[0-9]{2}[\:]{1}[0-9]{2}([\.]{1}[0-9]+)?[\+]{1}[0-9]{2}[\:]{1}[0-9]{2}$`
	// timeZonePattern    = `^[a-zA-Z]+/[a-z\-\_+\-A-Z]+$`
	timeZonePattern = `^[a-zA-Z0-9\-−_\/\+]+$`
)

var (
	// mailRegexp        = regexp.MustCompile(mailPattern)
	dateRegexp             = regexp.MustCompile(datePattern)
	dateTimeRegexp         = regexp.MustCompile(dateTimePattern)
	timeWithLocationRegexp = regexp.MustCompile(timeWithLocationPattern)
	timeZoneRegexp         = regexp.MustCompile(timeZonePattern)
)

// ValidateDatetimeType validate if the value is a datetime type
func ValidateDatetimeType(value interface{}) error {
	// time type is supported
	if _, ok := value.(time.Time); ok {
		return nil
	}

	// timestamp type is supported
	if IsNumeric(value) {
		return nil
	}

	// string type with time format is supported
	if _, ok := IsTime(value); ok {
		return nil
	}
	return fmt.Errorf("value(%+v) is not of time type", value)
}

func IsDate(sInput interface{}) bool {
	switch val := sInput.(type) {
	case string:
		if len(val) == 0 {
			return false
		}
		return dateRegexp.MatchString(val)
	default:
		return false
	}
}

// IsTime 是否是时间类型
func IsTime(sInput interface{}) (DateTimeFieldType, bool) {
	switch val := sInput.(type) {
	case string:
		if dateTimeRegexp.MatchString(val) {
			return TimeWithoutLocationType, true
		}
		if timeWithLocationRegexp.MatchString(val) {
			return TimeWithLocationType, true
		}
		return InvalidDateTimeType, false
	default:
		return InvalidDateTimeType, false
	}
}

// IsTimeZone 是否是时区类型
func IsTimeZone(sInput interface{}) bool {
	switch val := sInput.(type) {
	case string:
		if len(val) == 0 {
			return false
		}
		return timeZoneRegexp.MatchString(val)
	default:
		return false
	}
}
