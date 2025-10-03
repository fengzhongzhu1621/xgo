package zaplogger

import (
	"testing"
	"time"

	"github.com/gookit/goutil/testutil/assert"
)

func TestCustomTimeFormat(t *testing.T) {
	date := time.Date(2006, 1, 2, 15, 4, 5, 0, time.Local)
	dateStr := CustomTimeFormat(date, "2006-01-02 15:04:05.000")
	assert.Equal(t, dateStr, "2006-01-02 15:04:05.000")
}

func TestDefaultTimeFormat(t *testing.T) {
	date := time.Date(2006, 1, 2, 15, 4, 5, 0, time.Local)
	dateStr := string(DefaultTimeFormat(date))
	assert.Equal(t, dateStr, "2006-01-02 15:04:05.000")
}

func BenchmarkDefaultTimeFormat(b *testing.B) {
	t := time.Now()
	for i := 0; i < b.N; i++ {
		DefaultTimeFormat(t)
	}
}

func BenchmarkCustomTimeFormat(b *testing.B) {
	t := time.Now()
	for i := 0; i < b.N; i++ {
		CustomTimeFormat(t, "2006-01-02 15:04:05.000")
	}
}

func TestGetLogEncoderKey(t *testing.T) {
	tests := []struct {
		name   string
		defKey string
		key    string
		want   string
	}{
		{"custom", "T", "Time", "Time"},
		{"default", "T", "", "T"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLogEncoderKey(tt.defKey, tt.key); got != tt.want {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestNewTimeEncoder(t *testing.T) {
	encoder := NewTimeEncoder("")
	assert.NotNil(t, encoder)

	encoder = NewTimeEncoder("2006-01-02 15:04:05")
	assert.NotNil(t, encoder)

	tests := []struct {
		name string
		fmt  string
	}{
		{"seconds timestamp", "seconds"},
		{"milliseconds timestamp", "milliseconds"},
		{"nanoseconds timestamp", "nanoseconds"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTimeEncoder(tt.fmt)
			assert.NotNil(t, got)
		})
	}
}
