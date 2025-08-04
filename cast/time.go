package cast

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/fengzhongzhu1621/xgo/validator"
)

type timeFormatType int

const (
	timeFormatNoTimezone timeFormatType = iota
	timeFormatNamedTimezone
	timeFormatNumericTimezone
	timeFormatNumericAndNamedTimezone
	timeFormatTimeOnly
)

const (
	dateTimePattern         = `^[0-9]{4}[\-]{1}[0-9]{2}[\-]{1}[0-9]{2}[\s]{1}[0-9]{2}[\:]{1}[0-9]{2}[\:]{1}[0-9]{2}$`
	timeWithLocationPattern = `^[0-9]{4}[\-]{1}[0-9]{2}[\-]{1}[0-9]{2}[T]{1}[0-9]{2}[\:]{1}[0-9]{2}[\:]{1}[0-9]{2}([\.]{1}[0-9]+)?[\+]{1}[0-9]{2}[\:]{1}[0-9]{2}$`
)

var (
	dateTimeRegexp         = regexp.MustCompile(dateTimePattern)
	timeWithLocationRegexp = regexp.MustCompile(timeWithLocationPattern)
)

var timeFormats = []timeFormat{
	{time.RFC3339, timeFormatNumericTimezone},
	{"2006-01-02T15:04:05", timeFormatNoTimezone}, // iso8601 without timezone
	{time.RFC1123Z, timeFormatNumericTimezone},
	{time.RFC1123, timeFormatNamedTimezone},
	{time.RFC822Z, timeFormatNumericTimezone},
	{time.RFC822, timeFormatNamedTimezone},
	{time.RFC850, timeFormatNamedTimezone},
	{"2006-01-02 15:04:05.999999999 -0700 MST", timeFormatNumericAndNamedTimezone}, // Time.String()
	{
		"2006-01-02T15:04:05-0700",
		timeFormatNumericTimezone,
	}, // RFC3339 without timezone hh:mm colon
	{
		"2006-01-02 15:04:05Z0700",
		timeFormatNumericTimezone,
	}, // RFC3339 without T or timezone hh:mm colon
	{"2006-01-02 15:04:05", timeFormatNoTimezone},
	{time.ANSIC, timeFormatNoTimezone},
	{time.UnixDate, timeFormatNamedTimezone},
	{time.RubyDate, timeFormatNumericTimezone},
	{"2006-01-02 15:04:05Z07:00", timeFormatNumericTimezone},
	{"2006-01-02", timeFormatNoTimezone},
	{"02 Jan 2006", timeFormatNoTimezone},
	{"2006-01-02 15:04:05 -07:00", timeFormatNumericTimezone},
	{"2006-01-02 15:04:05 -0700", timeFormatNumericTimezone},
	{time.Kitchen, timeFormatTimeOnly},
	{time.Stamp, timeFormatTimeOnly},
	{time.StampMilli, timeFormatTimeOnly},
	{time.StampMicro, timeFormatTimeOnly},
	{time.StampNano, timeFormatTimeOnly},
}

type timeFormat struct {
	format string
	typ    timeFormatType
}

func (f timeFormat) hasTimezone() bool {
	// We don't include the formats with only named timezones, see
	// https://github.com/golang/go/issues/19694#issuecomment-289103522
	return f.typ >= timeFormatNumericTimezone && f.typ <= timeFormatNumericAndNamedTimezone
}

// ToTimeE casts an interface to a time.Time type.
// 时间格式转换，转换为UTC时区.
func ToTimeE(i interface{}) (tim time.Time, err error) {
	return ToTimeInDefaultLocationE(i, time.UTC)
}

// ToTime casts an interface to a time.Time type.
func ToTime(i interface{}) time.Time {
	v, _ := ToTimeE(i)
	return v
}

// ToTimeInDefaultLocationE casts an empty interface to time.Time,
// interpreting inputs without a timezone to be in the given location,
// or the local timezone if nil.
func ToTimeInDefaultLocationE(i interface{}, location *time.Location) (tim time.Time, err error) {
	i = indirect(i)

	switch v := i.(type) {
	case time.Time:
		return v, nil
	case string:
		return StringToDateInDefaultLocation(v, location)
	case int:
		return time.Unix(int64(v), 0), nil
	case int64:
		return time.Unix(v, 0), nil
	case int32:
		return time.Unix(int64(v), 0), nil
	case uint:
		return time.Unix(int64(v), 0), nil
	case uint64:
		return time.Unix(int64(v), 0), nil
	case uint32:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("unable to cast %#v of type %T to Time", i, i)
	}
}

// 时间差格式转换
// ToDurationE casts an interface to a time.Duration type.
func ToDurationE(i interface{}) (d time.Duration, err error) {
	i = indirect(i)

	switch s := i.(type) {
	case time.Duration:
		return s, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		d = time.Duration(ToInt64(s))
		return
	case float32, float64:
		d = time.Duration(ToFloat64(s))
		return
	case string:
		// 判断字符串 s 中是否包含 chars 中的任何一个字符
		if strings.ContainsAny(s, "nsuµmh") {
			d, err = time.ParseDuration(s)
		} else {
			d, err = time.ParseDuration(s + "ns")
		}
		return
	default:
		err = fmt.Errorf("unable to cast %#v of type %T to Duration", i, i)
		return
	}
}

func ToTimeInDefaultLocation(i interface{}, location *time.Location) time.Time {
	v, _ := ToTimeInDefaultLocationE(i, location)
	return v
}

// ToDuration casts an interface to a time.Duration type.
func ToDuration(i interface{}) time.Duration {
	v, _ := ToDurationE(i)
	return v
}

// ToDurationSlice casts an interface to a []time.Duration type.
func ToDurationSlice(i interface{}) []time.Duration {
	v, _ := ToDurationSliceE(i)
	return v
}

// ToDurationSliceE casts an interface to a []time.Duration type.
func ToDurationSliceE(i interface{}) ([]time.Duration, error) {
	if i == nil {
		return []time.Duration{}, fmt.Errorf(
			"unable to cast %#v of type %T to []time.Duration",
			i,
			i,
		)
	}

	switch v := i.(type) {
	case []time.Duration:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]time.Duration, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := ToDurationE(s.Index(j).Interface())
			if err != nil {
				return []time.Duration{}, fmt.Errorf(
					"unable to cast %#v of type %T to []time.Duration",
					i,
					i,
				)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []time.Duration{}, fmt.Errorf(
			"unable to cast %#v of type %T to []time.Duration",
			i,
			i,
		)
	}
}

// StringToDate attempts to parse a string into a time.Time type using a
// predefined list of formats.  If no suitable format is found, an error is
// returned.
func StringToDate(s string) (time.Time, error) {
	return parseDateWith(s, time.UTC, timeFormats)
}

// 将字符串格式的时间转换为 time.Time 类型
// StringToDateInDefaultLocation casts an empty interface to a time.Time,
// interpreting inputs without a timezone to be in the given location,
// or the local timezone if nil.
func StringToDateInDefaultLocation(s string, location *time.Location) (time.Time, error) {
	return parseDateWith(s, location, timeFormats)
}

// parseDateWith 将字符串格式的时间转换为 time.Time 类型，根据 timeFormats 的格式进行转换.
func parseDateWith(s string, location *time.Location, formats []timeFormat) (d time.Time, e error) {
	for _, format := range formats {
		if d, e = time.Parse(format.format, s); e == nil {
			// Some time formats have a zone name, but no offset, so it gets
			// put in that zone name (not the default one passed in to us), but
			// without that zone's offset. So set the location manually.
			if format.typ <= timeFormatNamedTimezone {
				if location == nil {
					location = time.Local
				}
				year, month, day := d.Date()
				hour, min, sec := d.Clock()
				d = time.Date(year, month, day, hour, min, sec, d.Nanosecond(), location)
			}

			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}

// ConvToTime convert value to time type
func ConvToTime(value interface{}) (time.Time, error) {
	timeVal, ok := value.(time.Time)
	if ok {
		return timeVal, nil
	}

	if validator.IsNumeric(value) {
		intVal, err := ToInt64E(value)
		if err != nil {
			return time.Time{}, fmt.Errorf("value(%+v) type is not int", value)
		}
		return time.Unix(intVal, 0), nil
	}

	strVal, ok := value.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("value(%+v) type is not string", value)
	}

	timeType, isTime := isTime(strVal)
	if !isTime {
		return time.Time{}, fmt.Errorf("value(%+v)  is not a string of time type", value)
	}

	return Str2Time(strVal, timeType), nil
}

// Str2Time string convert to time type
func Str2Time(timeStr string, timeType validator.DateTimeFieldType) time.Time {
	var layout string
	switch timeType {
	case validator.TimeWithoutLocationType:
		layout = "2006-01-02 15:04:05"
	case validator.TimeWithLocationType:
		layout = "2006-01-02T15:04:05+08:00"
	default:
		return time.Time{}
	}

	fTime, err := time.ParseInLocation(layout, timeStr, time.Local)
	if nil != err {
		return fTime
	}
	return fTime.UTC()
}

func isTime(sInput interface{}) (validator.DateTimeFieldType, bool) {
	switch val := sInput.(type) {
	case string:
		if dateTimeRegexp.MatchString(val) {
			return validator.TimeWithoutLocationType, true
		}
		if timeWithLocationRegexp.MatchString(val) {
			return validator.TimeWithLocationType, true
		}
		return validator.InvalidDateTimeType, false
	default:
		return validator.InvalidDateTimeType, false
	}
}
