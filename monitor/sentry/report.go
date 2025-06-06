package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
)

var sentryOn bool

// SetSentryCapture init the sentryEnabled var
func SetSentryCapture(value bool) {
	sentryOn = value
}

func EnableSentry() {
	sentryOn = true
}

func DisableSentry() {
	sentryOn = false
}

// ReportToSentry is a shortcut to build and send an event to sentry
func ReportToSentry(message string, extra map[string]interface{}) {
	event := sentry.NewEvent()
	event.Message = message
	event.Level = "error"
	event.Timestamp = time.Now()
	event.Extra = extra
	if sentryOn {
		sentry.CaptureEvent(event)
	}
}
